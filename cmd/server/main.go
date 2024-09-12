package main

import (
    "github.com/labstack/echo/v4"
    "p2pbot/internal/app"
    "p2pbot/internal/handlers"
    echojwt "github.com/labstack/echo-jwt/v4"
    "p2pbot/internal/JWTConfig"
    "p2pbot/internal/db/repository"
    "p2pbot/internal/services"
    "p2pbot/internal/utils"
)


func main() {
    DB, cfg, err := app.Init()
    if err != nil {
        panic(err)
    }

    userRepo := repository.NewUserRepository(DB)
    userService := services.NewUserService(userRepo)
    trackerRepo := repository.NewTrackerRepository(DB)
    trackerService := services.NewTrackerService(trackerRepo)

	binance := services.NewBinanceExchange(cfg)
	bybit := services.NewBybitExcahnge(cfg)

    controller := handlers.NewController(userService,
                                            trackerService,
                                            map[string]services.ExchangeI{
                                                "binance": binance,
                                                "bybit": bybit,
                                            },
                                            cfg.Website.JWTSecret)


    utils.NewLogger()
    e := echo.New()
    e.Use(utils.LoggingMiddleware)

    publicGroup := e.Group("/api/v1/public")

    publicGroup.POST("/login", controller.Login) 
    publicGroup.POST("/signup", controller.Signup)

    privateGroup := e.Group("/api/v1/private")

    config := JWTConfig.NewJWTConfig(cfg)
    privateGroup.Use(echojwt.WithConfig(config))
    privateGroup.Use(utils.AuthMiddleware)

    privateGroup.GET("/trackers", controller.GetTrackers)
    privateGroup.POST("/trackers", controller.CreateTracker)
    

    e.Logger.Fatal(e.Start(":1323"))
}
