package JWTConfig

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"p2pbot/internal/config"
)

type JWTCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTConfig(cfg *config.Config) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTCustomClaims)
		},
		SigningKey:  []byte(cfg.Website.JWTSecret),
		TokenLookup: "cookie:token",
	}
}
