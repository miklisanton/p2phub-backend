package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"p2pbot/internal/bot"
	"p2pbot/internal/config"
	"p2pbot/internal/db/drivers"
	"p2pbot/internal/db/repository"
	"p2pbot/internal/services"
	"time"
)

var (
	DB  *sqlx.DB
	cfg *config.Config
)

func init() {

	cfg, err := config.NewConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	connectionURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host+":"+cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSL)
	DB, err = drivers.Connect(connectionURL)
	if err != nil {
		log.Fatal(err)
	}

}

// Delete keyboard message after send
func main() {

	//url := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	//payload := `{"fiat":"CZK","page":1,"rows":10,"tradeType":"BUY","asset":"USDT","countries":[],"proMerchantAds":false,"shieldMerchantAds":false,"filterType":"all","periods":[],"additionalKycVerifyFilter":0,"publisherType":null,"payTypes":[],"classifies":["mass","profession"]}`

	trackerRepo := repository.NewTrackerRepository(DB)
	userRepo := repository.NewUserRepository(DB)

	trackerService := services.NewTrackerService(trackerRepo)
	userService := services.NewUserService(userRepo)

	cli, err := config.ParseCLI()
	if err != nil {
		log.Fatal("Error parsing cli: ", err)
	}

	cfg, err := config.NewConfig(cli)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	//Supported exchanges
	binance := services.NewBinanceExchange(cfg)
	bybit := services.NewBybitExcahnge(cfg)

	exs := []services.ExchangeI{binance, bybit}

	tgbot, err := bot.NewBot(cfg, userService, trackerService, exs)
	if err != nil {
		log.Fatal("Error starting bot: ", err)
	}

	go func() {
		err = tgbot.MonitorAds(time.Minute * 1)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		tgbot.NotifyUsers()
	}()
	tgbot.Start()

}
