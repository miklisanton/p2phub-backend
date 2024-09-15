package handlers

import (
	"p2pbot/internal/services"
)

type Controller struct {
    userService *services.UserService
    trackerService *services.TrackerService
    exchanges map[string]services.ExchangeI
    JWTSecret string
    TgLink string
}

func NewController(userService *services.UserService,
                    trackerService *services.TrackerService,
                    exs map[string]services.ExchangeI,
                    JWTSecret string,
                    tgLink string) *Controller {
    return &Controller{userService, trackerService, exs, JWTSecret, tgLink}
}
