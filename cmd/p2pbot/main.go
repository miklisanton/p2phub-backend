package main

import (
    "log"
    "p2pbot/internal/app"
    "p2pbot/internal/bot"
    "p2pbot/internal/db/repository"
    "p2pbot/internal/rediscl"
    "p2pbot/internal/services"
    "p2pbot/internal/rabbitmq"
    "p2pbot/internal/utils"
    "time"
)

// Delete keyboard message after send
func main() {
    utils.NewLogger()

    // wait until all services are up
    time.Sleep(10 * time.Second)
    DB, cfg, err := app.Init()
    if err != nil {
        panic(err)
    }

    //url := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
    //payload := `{"fiat":"CZK","page":1,"rows":10,"tradeType":"BUY","asset":"USDT","countries":[],"proMerchantAds":false,"shieldMerchantAds":false,"filterType":"all","periods":[],"additionalKycVerifyFilter":0,"publisherType":null,"payTypes":[],"classifies":["mass","profession"]}`

    trackerRepo := repository.NewTrackerRepository(DB)
    userRepo := repository.NewUserRepository(DB)

    trackerService := services.NewTrackerService(trackerRepo)
    userService := services.NewUserService(userRepo)

    rediscl.InitRedisClient(cfg.Redis.Host, cfg.Redis.Port)

    //Supported exchanges
    binance := services.NewBinanceExchange(cfg)
    bybit := services.NewBybitExcahnge(cfg)
    exs := []services.ExchangeI{binance, bybit}

    tgbot, err := bot.NewBot(cfg, userService, trackerService, exs)
    if err != nil {
        log.Fatal("Error starting bot: ", err)
    }
    // Rabbitmq setup
    rabbit, err := rabbitmq.NewRabbitMQ(cfg)
    if err != nil {
        log.Fatal("Error starting rabbitmq: ", err)
    }
    if err := rabbit.DeclareExchange("notifications"); err != nil {
        log.Fatal("Error declaring exchange: ", err)
    }
    queueName, err := rabbit.QueueBindNDeclare()
    if err != nil {
        log.Fatal("Error declaring queue: ", err)
    }
    rabbit.StartConsuming(queueName, tgbot.HandleNotification)

    tgbot.Start()

}
