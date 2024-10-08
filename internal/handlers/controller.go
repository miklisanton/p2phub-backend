package handlers

import (
	"p2pbot/internal/services"
    "p2pbot/internal/config"
)

type Controller struct {
    userService *services.UserService
    trackerService *services.TrackerService
    exchanges map[string]services.ExchangeI
    JWTSecret string
    TgLink string
    SubPrice string
    SubCurrency string
}

func NewController(userService *services.UserService,
                    trackerService *services.TrackerService,
                    exs map[string]services.ExchangeI,
                    cfg *config.Config) *Controller {

    return &Controller{
        userService, 
        trackerService,
        exs,
        cfg.Website.JWTSecret, 
        cfg.Telegram.InviteLink,
        cfg.Website.SubPrice,
        cfg.Website.SubCurrency,
    }

}
