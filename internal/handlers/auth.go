package handlers

import (
	"database/sql"
    "p2pbot/internal/rediscl"
	"net/http"
	"p2pbot/internal/JWTConfig"
	"p2pbot/internal/db/models"
	"p2pbot/internal/requests"
	"p2pbot/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/teris-io/shortid"
)


func (contr *Controller) Signup(c echo.Context) error {
    u := new(requests.LoginRequest)
    if err := c.Bind(u); err != nil {
        return err
    }

    v := validator.New()
    err := v.Struct(u)
    if err != nil {
        out := make(map[string]string)
        for _, e := range err.(validator.ValidationErrors) {
            out[e.Field()] = e.Tag() + " " + e.Param()
        }    
        return c.JSON(http.StatusBadRequest, echo.Map{
            "message": "Validation error",
            "errors": out,
        })
    }
    
    _, err = contr.userService.GetUserByEmail(u.Email)

    if err == nil {
        return c.JSON(http.StatusConflict, map[string]any{
            "message": "User with this email already exists",
            "errors": map[string]any{
                "email": "email taken",
            },
        },
    )}

    passwordHash, err := utils.HashPassword(u.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]any{
            "message": "Internal server error",
            "errors": map[string]any{
                "internal": "failed to create user",
            },
        },
    )}

    _, err = contr.userService.CreateUser(&models.User{
        Email: &u.Email,
        Password_en: &passwordHash,
    })

    if err != nil {
        return err
    }

    return c.JSON(http.StatusCreated, map[string]any{
        "message": "User created, login",
        "user": map[string]any{
            "email": u.Email,
        },
    })
}

func (contr *Controller) Login(c echo.Context) error {
    u := new(requests.LoginRequest)
    if err := c.Bind(u); err != nil {
        return err
    }

    user, err := contr.userService.GetUserByEmail(u.Email)

    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]any{
            "message": "Invalid email or password",
            "errors": map[string]any{
                "credentials": "Invalid email or password",
            },
        },
    )}

    if ok := utils.CheckPasswordHash(u.Password, *user.Password_en); !ok {
        return c.JSON(http.StatusUnauthorized, map[string]any{
            "message": "Invalid email or password",
            "errors": map[string]any{
                "credentials": "Invalid email or password",
            },
        },
    )}

    claims := &JWTConfig.JWTCustomClaims{
        Email: *user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(contr.JWTSecret))
    if err != nil {
        return err
    }

    c.SetCookie(&http.Cookie{
        HttpOnly: true,
        Value: tokenString,
        Expires: time.Now().Add(time.Hour * 24),
        SameSite: http.SameSiteStrictMode,
        Name: "token",
    })

    return c.JSON(http.StatusOK, map[string]any{
        "message": "Login successful",
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
        "link": link,
    })         
}

func (cont *Controller) GetCSRFToken(c echo.Context) error {
    csrf := c.Get("csrf").(string)
    return c.JSON(http.StatusOK, map[string]any{
        "csrf": csrf,
    })
}
