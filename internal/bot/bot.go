package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"p2pbot/internal/config"
	"p2pbot/internal/db/models"
	"p2pbot/internal/fsm"
	"p2pbot/internal/services"
	"p2pbot/internal/utils"
	"strings"
	"time"
)

type Bot struct {
	api            *tgbotapi.BotAPI
	userService    *services.UserService
	trackerService *services.TrackerService
	NotificationCh chan utils.Notification
	exchanges      []services.ExchangeI
	toDelete       []int
}

func NewBot(cfg *config.Config, userSvc *services.UserService, trackerSvc *services.TrackerService, exs []services.ExchangeI) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Telegram.APIkey)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:            api,
		userService:    userSvc,
		trackerService: trackerSvc,
		NotificationCh: make(chan utils.Notification),
		exchanges:      exs}, nil
}

func (bot *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.api.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	sm := fsm.New()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID

		switch update.Message.Command() {
		case "":
			//handle user input
			msg := update.Message.Text
			bot.toDelete = append(bot.toDelete, update.Message.MessageID)
			log.Println(msg)

			tracker := bot.trackerService.GetTrackerStaging(chatID)

			var nextState fsm.State

			currentState := sm.GetState(chatID)
			switch currentState {
			case fsm.Welcome:
				id := bot.SendMessage(chatID, "/new to create new tracker")
				bot.toDelete = append(bot.toDelete, id)
			case fsm.AwaitingExchange:
				if _, ok := bot.trackerService.Exchanges[msg]; ok {
					nextState, err = sm.Transition(chatID, fsm.ExchangeFound, chatID, msg, tracker)
					if err != nil {
						log.Fatal(err)
					}

					id := bot.SendMessage(chatID, fmt.Sprintf("Provide currency symbol, <EUR> for example"))
					bot.toDelete = append(bot.toDelete, id)
				} else {
					id := bot.SendMessage(chatID, fmt.Sprintf("Exhcange %s not supported", msg))
					bot.toDelete = append(bot.toDelete, id)

					nextState, err = sm.Transition(chatID, fsm.ExchangeNotFound)
					if err != nil {
						log.Fatal(err)
					}
				}
			case fsm.AwaitingСurrency:
				nextState, err = sm.Transition(chatID, fsm.CurrencyGiven, strings.ToUpper(msg), tracker)
				if err != nil {
					log.Fatal(err)
				}

				id := bot.SendMessage(chatID, fmt.Sprintf("Provide username on %s", tracker.Exchange))
				bot.toDelete = append(bot.toDelete, id)
			case fsm.AwaitingExchangeUsername:
				nextState, err = sm.Transition(chatID, fsm.UsernameGiven)
				if err != nil {
					log.Fatal(err)
				}

				bybitName, binanceName := "", ""
				if tracker.Exchange == "bybit" {
					bybitName = msg
				}
				if tracker.Exchange == "binance" {
					binanceName = msg
				}

				err = bot.userService.CreateUser(&models.User{
					ChatID:      chatID,
					BinanceName: binanceName,
					BybitName:   bybitName,
				})
				if err != nil {
					log.Fatal(err)
				}

				id := bot.SendMessage(chatID, fmt.Sprintf("Provide side BUY/SELL"))
				bot.toDelete = append(bot.toDelete, id)
			case fsm.AwaitingSide:
				tracker.Side = strings.ToUpper(msg)
				for _, exchange := range bot.exchanges {
					if strings.ToLower(exchange.GetName()) == tracker.Exchange {
						user, err := bot.userService.GetUserById(chatID)
						if err != nil {
							log.Fatal(err)
						}

						username, err := utils.GetField(user, exchange.GetName()+"Name")
						if err != nil {
							log.Fatalf("error getting exchange name: %s", err)
						}

						adv, err := exchange.GetAdvByName(tracker.Currency, tracker.Side, username.(string))
						if err != nil {
							id := bot.SendMessage(chatID, err.Error())
							bot.toDelete = append(bot.toDelete, id)

							nextState, err = sm.Transition(chatID, fsm.AdvertisementNotFound)
							if err != nil {
								log.Fatal(err)
							}

							id = bot.SendMessage(chatID, "/new to create new tracker")
							bot.toDelete = append(bot.toDelete, id)
						} else {
							tracker.Payment = adv.GetPaymentMethods()

							err = bot.trackerService.CreateTracker(tracker)
							if err != nil {
								log.Fatal(err)
							}
							nextState, err = sm.Transition(chatID, fsm.AdvertisementFound)
							if err != nil {
								log.Fatal(err)
							}

							if err := bot.DeleteMessages(chatID); err != nil {
								log.Fatal(err)
							}
							bot.SendMessage(chatID, "Tracker created successfully")
						}
					}
				}

			default:
				log.Fatalf("not implemented %d", currentState)
			}

			log.Println("next state:", nextState)
		case "new":
			if sm.GetState(chatID) == fsm.Welcome {
				id := bot.SendMessage(chatID, "Choose exchange bybit/binance")
				bot.toDelete = append(bot.toDelete, id)

				s, err := sm.Transition(chatID, fsm.NewTracker)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("next state:", s)
			}
		case "start":
			bot.SendMessage(chatID, "/new to create new tracker")
		default:
			bot.SendMessage(chatID, "Unknown command")
		}
	}
}

func (bot *Bot) SendMessage(chatID int64, text string) int {
	msg := tgbotapi.NewMessage(chatID, text)
	msgSent, err := bot.api.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return 0
	}

	return msgSent.MessageID
}

func (bot *Bot) SendMultiple(ids []int64, text string) {
	for _, id := range ids {
		bot.SendMessage(id, text)
	}
}

func (bot *Bot) MonitorAds(refreshRate time.Duration) error {
	ticker := time.NewTicker(refreshRate)
	for range ticker.C {
		trackers, err := bot.trackerService.GetAllTrackers()
		if err != nil {
			return fmt.Errorf("error getting trackers: %s", err)
		}
		for _, tracker := range trackers {
			for _, exchange := range bot.exchanges {
				if strings.ToLower(exchange.GetName()) != tracker.Exchange {
					continue
				}
				// Retrieve username on exchange
				username, err := utils.GetField(tracker, exchange.GetName()+"Name")
				if err != nil {
					return fmt.Errorf("error getting exchange name: %s", err)
				}
				// Skip, username not provided
				if username == "" {
					continue
				}

				exchangeBestAdv, err := exchange.GetBestAdv(tracker.Currency, tracker.Side, tracker.Payment)
				if err != nil {
					return err
				}
				// Send notification if best advertisement doesn't match with tracker username
				if exchangeBestAdv.GetName() != username {
					// Notify if user reacted on previous notification(updated his adv)
					if !tracker.Waiting {
						bot.NotificationCh <- utils.Notification{
							ChatID:    tracker.ChatID,
							Data:      exchangeBestAdv,
							Exchange:  tracker.Exchange,
							Direction: tracker.Side,
							Currency:  tracker.Currency,
						}
					}
					// Set waiting_adv flag, until user puts his advertisement on top
					if err := bot.trackerService.SetWaitingFlag(tracker.ID, true); err != nil {
						return err
					}
				} else {
					// User updated his adv, enable notifications
					if err := bot.trackerService.SetWaitingFlag(tracker.ID, false); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (bot *Bot) NotifyUsers() {
	for notification := range bot.NotificationCh {
		q, minAmount, maxAmount := notification.Data.GetQuantity()
		c := notification.Currency
		msg := fmt.Sprintf("Your %s order on %s was outbided by %s.\nPrice: %.3f %s\nQuantity: %.2fUSDT\nMin: %.2f%s\nMax: %.2f%s",
			notification.Direction, notification.Exchange, notification.Data.GetName(), notification.Data.GetPrice(), c, q, minAmount, c, maxAmount, c)

		bot.SendMessage(notification.ChatID, msg)
		bot.SendMessage(notification.ChatID, "You won't get notifications until you update your order")
	}
}

func (bot *Bot) DeleteMessages(chatID int64) error {
	for _, msgID := range bot.toDelete {
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
		_, err := bot.api.Send(deleteMsg)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}

	bot.toDelete = []int{}
	return nil
}
