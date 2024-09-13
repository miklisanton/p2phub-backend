package handlers

import (
    "p2pbot/internal/utils"
    "fmt"
	"net/http"
	"github.com/labstack/echo/v4"
    "p2pbot/internal/requests"
    "p2pbot/internal/db/models"
)

func (contr *Controller) GetTrackers(c echo.Context) error {
    email := c.Get("email").(string)
    u, err := contr.userService.GetUserByEmail(email)
    if err != nil {
        return err
    }

    trackers, err := contr.trackerService.GetTrackersByUserId(u.ID)
    if err != nil {
        return err
    }
    return c.JSON(http.StatusOK, map[string]any{
        "message": fmt.Sprintf("Trackers for user %s", email),
        "trackers": trackers,  
    })
}

func (contr *Controller) CreateTracker(c echo.Context) error {
    email := c.Get("email").(string)
    u, err := contr.userService.GetUserByEmail(email)
    if err != nil {
        return err
    }

    trackerReq := new(requests.TrackerRequest)
    if err := c.Bind(trackerReq); err != nil {
        return err
    }

    // Create tracker entity and validate its fields
    tracker := &models.Tracker{
        UserID: u.ID,
        Exchange: trackerReq.Exchange,
        Currency: trackerReq.Currency,
        Side: trackerReq.Side,
        Username: trackerReq.Username,
    }

    if err := contr.trackerService.ValidateTracker(tracker, false); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]any{
            "message": "Validation error",
            "errors": map[string]any{
                "invalid_param": err.Error(),
            },
        })
    }
    // Check if exchange is supported
    exchange, ok := contr.exchanges[tracker.Exchange]
    if !ok {
        return c.JSON(http.StatusBadRequest, map[string]any{
            "message": "exchange not found",
            "errors": map[string]any{
                "exchange": fmt.Sprintf("%s not supported", tracker.Exchange),
            },
        })
    }
    
    ads, err := exchange.GetAdsByName(tracker.Currency,
                                        tracker.Side,
                                        tracker.Username)
    if err != nil {
        return c.JSON(http.StatusNotFound, map[string]any{
            "message": "error getting ads",
            "errors": map[string]any{
                exchange.GetName(): err.Error(),
            },
        })
    }

    utils.Logger.LogInfo().Fields(map[string]interface{}{
        "email": email,
        exchange.GetName(): ads,
    }).Msg("Ads found")

    createdTrackers := make([]*models.Tracker, 0)
    for _, adv := range ads {
        tracker.Payment = adv.GetPaymentMethods()

        err = contr.trackerService.CreateTracker(tracker)
        createdTrackers = append(createdTrackers, tracker)
        if err != nil {
            return err
        }
    }

    return c.JSON(http.StatusCreated, map[string]any{
        "message": "Trackers created",
        "trackers": createdTrackers,
    })
}

func (contr *Controller) GetPaymentMethods (c echo.Context) error {
    email := c.Get("email").(string)
    // Check query parameters
    exchange := c.QueryParam("exchange")
    if exchange == "" {
        return c.JSON(http.StatusBadRequest, map[string]any{
            "message": "exchange not found",
            "errors": map[string]any{
                "exchange": "query parameter not provided",
            },
                
        })
    }
    exch, ok := contr.exchanges[exchange]
    if !ok {
        return c.JSON(http.StatusBadRequest, map[string]any{
            "message": "exchange not found",
            "errors": map[string]any{
                "exchange": fmt.Sprintf("%s not supported", exchange),
            },
        })
    }

    currency := c.QueryParam("currency")
    if currency == "" {
        return c.JSON(http.StatusBadRequest, map[string]any{
            "message": "currency not found",
            "errors": map[string]any{
                "currency": "query parameter not provided",
            },
        })
    }

    supportedCurrencies, err := exch.GetCachedCurrencies()
    if err != nil {
        return err
    }
    if !utils.Contains(supportedCurrencies, currency) {
        return c.JSON(http.StatusBadRequest, map[string]any{
            "message": "currency not found",
            "errors": map[string]any{
                "currency": fmt.Sprintf("%s not supported", currency),
            },
        })
    }
    out, err := exch.GetCachedPaymentMethods(currency)
    if err != nil {
        return err
    }

    utils.Logger.LogInfo().Fields(map[string]interface{}{
        "email": email,
        "exchange": exchange,
        "currency": currency,
        "options": out,
    }).Msg("Payment methods requested")

    return c.JSON(http.StatusOK, map[string]any{
        "message": "Form options",
        "options": out,
    })
}

func (contr *Controller) GetCurrencies(c echo.Context) error {
    email := c.Get("email").(string)
    // Check query parameters
    exchange := c.QueryParam("exchange")
    if exchange == "" {
        return c.JSON(http.StatusBadRequest, map[string]any{
            "message": "exchange not found",
            "errors": map[string]any{
                "exchange": "query parameter not provided",
            },
                
        })
    }
    exch, ok := contr.exchanges[exchange]
    if !ok {
        return c.JSON(http.StatusBadRequest, map[string]any{
            "message": "exchange not found",
            "errors": map[string]any{
                "exchange": fmt.Sprintf("%s not supported", exchange),
            },
        })
    }

    out, err := exch.GetCachedCurrencies()
    if err != nil {
        return err
    }

    utils.Logger.LogInfo().Fields(map[string]interface{}{
        "email": email,
        "exchange": exchange,
        "options": out,
    }).Msg("Currencies requested")

    return c.JSON(http.StatusOK, map[string]any{
        "message": "Currencies",
        "options": out,
    })
}
