package main

import (
	"context"
	"fmt"
	"p2pbot/internal/app"
	"p2pbot/internal/db/repository"
	"p2pbot/internal/rabbitmq"
	"p2pbot/internal/services"
	"p2pbot/internal/tasks"
	"p2pbot/internal/utils"
	"time"
)

func main() {
    DB, cfg, err := app.Init()
    fmt.Println("DB: ", DB)
    if err != nil {
        panic(err)
    }
    utils.NewLogger()

    trackerRepo := repository.NewTrackerRepository(DB)
    userRepo := repository.NewUserRepository(DB)
    trackerService := services.NewTrackerService(trackerRepo)
    userService := services.NewUserService(userRepo)
    
    bybit := services.NewBybitExcahnge(cfg)
    binance := services.NewBinanceExchange(cfg)
    exs := []services.ExchangeI{binance, bybit}

    mapp, err := trackerService.GetIdsByCurrency("binance")
    if err != nil {
        fmt.Println("Error: ", err)
    }
    fmt.Println("map: ", mapp)
    rabbit, err := rabbitmq.NewRabbitMQ(cfg)
    if err != nil {
        fmt.Println("Error: ", err)
    }
    observer := tasks.NewAdsObserver(trackerService, userService, exs, rabbit)

    ctx := context.Background()
    observer.Start(1 * time.Minute, ctx)
}
