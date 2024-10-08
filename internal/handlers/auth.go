package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"p2pbot/internal/db/models"
	"p2pbot/internal/rediscl"
	"p2pbot/internal/requests"
	"p2pbot/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/teris-io/shortid"
)

// This POST request handler is used to sync auth0 users with database
func (contr *Controller) Signup(c echo.Context) error {
	u := new(requests.Auth0Request)
	if err := c.Bind(u); err != nil {
		return err
	}



    utils.Logger.Debug().Fields(map[string]interface{}{
        "email": u.Email,
        "secret": u.Secret,
    }).Msg("Signup request")

    if (u.Secret != os.Getenv("AUTH0_SIGNUP_SECRET")) {
        utils.Logger.LogError().Msg("Invalid secret in signup request")
        return c.JSON(http.StatusUnauthorized, map[string]any{
            "message": "Invalid secret",
        })
    }

    _, err := contr.userService.CreateUser(&models.User{
		Email:       &u.Email,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"message": "User created",
		"user": map[string]any{
			"email": u.Email,
		},
	})
}

func (contr *Controller) GetProfile(c echo.Context) error {
	email := c.Get("email").(string)
	u, err := contr.userService.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]any{
			"message": "User not found",
			"errors": map[string]any{
				"user": "not found",
			},
		})
	}
	if err != nil {
		return err
	}
    utils.Logger.LogInfo().Fields(map[string]interface{}{
        "email": u.Email,
        "chat_id": u.ChatID,
    }).Msg("User found")
	return c.JSON(http.StatusOK, map[string]any{
		"message": "User found",
		"user": map[string]any{
			"email":    u.Email,
			"telegram": u.ChatID,
		},
	})
}

func (contr *Controller) ConnectTelegram(c echo.Context) error {
	email := c.Get("email").(string)
	u, err := contr.userService.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]any{
			"message": "User not found",
			"errors": map[string]any{
				"user": "not found",
			},
		})
	}
	if err != nil {
		return err
	}
	// generate code with shortid
	code, err := shortid.Generate()
	if err != nil {
		return err
	}
	// save code to redis
	ctx := rediscl.RDB.Ctx
	if err := rediscl.RDB.Client.Set(ctx, "telegram_codes:"+code, u.ID, 15*time.Minute).Err(); err != nil {
		return err
	}
	// send link to user
	link := contr.TgLink + "?start=" + code
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Connect your telegram",
		"link":    link,
	})
}

func (cont *Controller) GetCSRFToken(c echo.Context) error {
	csrf := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, map[string]any{
		"csrf": csrf,
	})
}
