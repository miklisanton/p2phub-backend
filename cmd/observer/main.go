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
    // wait until all services are up
    time.Sleep(10 * time.Second)
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

    rabbit, err := rabbitmq.NewRabbitMQ(cfg)
    if err != nil {
        fmt.Println("Error: ", err)
    }
    observer := tasks.NewAdsObserver(trackerService, userService, exs, rabbit)

    ctx := context.Background()
    observer.Start(1 * time.Minute, ctx)
}
