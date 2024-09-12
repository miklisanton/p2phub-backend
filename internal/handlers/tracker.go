package handlers

import (
    "log"
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
    
    log.Printf("Find %d ads", len(ads))
    for _, adv := range ads {
        tracker.Payment = adv.GetPaymentMethods()

        err = contr.trackerService.CreateTracker(tracker)
        if err != nil {
            return err
        }
    }

    return c.String(http.StatusCreated, "TODO")
}
