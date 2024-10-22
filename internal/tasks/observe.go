package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"p2pbot/internal/db/models"
	"p2pbot/internal/rabbitmq"
	"p2pbot/internal/services"
	"p2pbot/internal/utils"
	"strings"
	"sync"
	"time"
)

type AdsObserver struct {
	trackerService *services.TrackerService
	userService    *services.UserService
	exchanges      []services.ExchangeI
	rabbitCl       *rabbitmq.RabbitMQ
}

func NewAdsObserver(trackerService *services.TrackerService, userService *services.UserService,
	exchanges []services.ExchangeI, rabbit *rabbitmq.RabbitMQ) *AdsObserver {
	return &AdsObserver{
		trackerService: trackerService,
		userService:    userService,
		exchanges:      exchanges,
		rabbitCl:       rabbit,
	}
}

func (ao *AdsObserver) Start(rate time.Duration, ctx context.Context) {
	if err := ao.rabbitCl.DeclareExchange("notifications"); err != nil {
		utils.Logger.LogError().Fields(map[string]interface{}{
			"error": err.Error(),
		}).Msg("Error declaring exchange")
	}
	ao.rabbitCl.Publish([]byte("Starting ads observer"))
	// Check ads with given rate
	ao.CheckAds()
	ticker := time.NewTicker(rate)
	for {
		select {
		case <-ticker.C:
			ao.CheckAds()
		case <-ctx.Done():
			return
		}
	}
}

func (ao *AdsObserver) CheckAds() {
	var wg sync.WaitGroup
	for _, ex := range ao.exchanges {

		utils.Logger.Debug().Msg(fmt.Sprintf("monitoring ads on %s", ex.GetName()))

		wg.Add(1)
		go func() {
			defer wg.Done()
			// Get map of "currency+side" -> [trackerID]
			idsMap, err := ao.trackerService.GetIdsByCurrency(strings.ToLower(ex.GetName()))
			utils.Logger.Debug().Fields(map[string]any{
				"map": idsMap,
			}).Msg("monitoring ads")
			if err != nil {
				return
			}
			ao.CheckAdsOnExchange(ex, idsMap)
		}()
	}
	wg.Wait()
}

func (ao *AdsObserver) CheckAdsOnExchange(ex services.ExchangeI, idsMap map[string][]int) {
	var wg sync.WaitGroup
	for key, ids := range idsMap {
		wg.Add(1)
		go func() error {
			defer wg.Done()
			// key, for example: "CZKSELL"
			currency := key[:3]
			side := key[3:]
			ads, err := ex.GetAds(currency, side)
			if err != nil {
				return err
			}
			for _, id := range ids {
				ao.CheckTracker(ads, id)
			}
			return nil
		}()
	}
	wg.Wait()
}

func (ao *AdsObserver) CheckTracker(ads []services.P2PItemI, trackerID int) {
	tracker, err := ao.trackerService.GetTrackerById(trackerID)
	if err != nil {
		return
	}
	if tracker.IsAggregated {
		for _, ad := range ads {
			if utils.ComparePaymentMethods(ad.GetPaymentMethods(), tracker.Payment) {
				// if advertisements payment methods contain one of the tracker payment methods
				if ad.GetName() != tracker.Username {
					// if advertisement name doesnt match tracker username
					// Notify user
					if !tracker.WaitingUpdate {
						ao.Notify(tracker, ad)
					}
					tracker.WaitingUpdate = true
					if err := ao.trackerService.CreateTracker(tracker); err != nil {
						log.Printf("Error updating tracker waiting update: %s", err)
					}
					return
				} else {
					// Tracked advertisement is the best advertisement across payment methods
					// Set outbidded to false
					tracker.WaitingUpdate = false
					tracker.Price = ad.GetPrice()
					log.Printf("User %s is not outbidded on %s", tracker.Username, tracker.Exchange)
					if err := ao.trackerService.CreateTracker(tracker); err != nil {
						log.Printf("Error updating tracker price: %s", err)
					}
					return
				}
			}
		}
	} else {
		for _, pMethod := range tracker.Payment {
			for _, ad := range ads {
				if utils.Contains(ad.GetPaymentMethods(), pMethod.Id) {
					if ad.GetName() != tracker.Username {
						//Notify user
						if !pMethod.Outbided {
							ao.Notify(tracker, ad)
						}
						//Set outbidded to true
						err = ao.trackerService.UpdateMethodOutbiddded(tracker.ID, pMethod.Id, true)
						if err != nil {
							log.Printf("Error updating outbidded status for %s on %s", pMethod.Id, tracker.Exchange)
						}
					} else {
						//set outbidded to false
						err := ao.trackerService.UpdateMethodOutbiddded(tracker.ID, pMethod.Id, false)
						if err != nil {
							log.Printf("Error updating outbidded status for %s on %s", pMethod.Id, tracker.Exchange)
						}
						log.Printf("User %s is not outbidded on %s for %s", tracker.Username, tracker.Exchange, pMethod.Id)
						//Update tracker price
						tracker.Price = ad.GetPrice()
						if err := ao.trackerService.CreateTracker(tracker); err != nil {
							log.Printf("Error updating tracker price: %s", err)
						} else {
							utils.Logger.LogDebug().Fields(map[string]interface{}{
								"id": tracker.ID,
							}).Msg("tracker updated")
						}
					}
					break
				}
			}
		}
	}
}

func (ao *AdsObserver) Notify(tracker *models.Tracker, ad services.P2PItemI) {
	user, err := ao.userService.GetUserByID(tracker.UserID)
	if err != nil {
		utils.Logger.LogError().Msg("Error retreiving user")
		return
	}
	if user.ChatID == nil {
		// if telegram  not connected
		utils.Logger.Info().Msg(fmt.Sprintf("Can't sent notification, because user with userID %d has no telegram connected", user.ID))
		return
	}
	if !tracker.Notify {
		return
	}
	n := utils.Notification{
		Data:     ad,
		Exchange: tracker.Exchange,
		Side:     tracker.Side,
		Currency: tracker.Currency,
		ChatID:   *user.ChatID,
	}
	nJson, err := json.Marshal(n)
	if err != nil {
		utils.Logger.LogError().Msg("Error converting user to json")
	}
	if err := ao.rabbitCl.Publish([]byte(nJson)); err != nil {
		utils.Logger.LogError().Fields(map[string]interface{}{
			"error": err.Error(),
		}).Msg("Error publishing message")
	}
}
