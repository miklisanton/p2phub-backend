package handlers

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"p2pbot/internal/rediscl"
	"p2pbot/internal/requests"
	"p2pbot/internal/utils"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/teris-io/shortid"
)

// Generatre order creates new invoice on cryptomus gateway.
// Unique order_id is stored in redis cache
// returns payment link
func (contr *Controller) CreateOrder(c echo.Context) error {
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
	// generate order id with shortid
	orderID, err := shortid.Generate()
	if err != nil {
		return err
	}
	utils.Logger.LogInfo().Fields(map[string]interface{}{
		"email":    u.Email,
		"order_id": orderID,
	}).Msg("Order ID generated")
	// save order id to redis
	ctx := rediscl.RDB.Ctx
	if err := rediscl.RDB.Client.Set(ctx, "order_id:"+orderID, u.ID, 2*time.Hour).Err(); err != nil {
		return err
	}
	// POST request to cryptomus gateway
	invoiceReq := requests.CryptomusRequest{
		OrderID:     orderID,
		Amount:      contr.SubPrice,
		Currency:    contr.SubCurrency,
		CallbackURL: os.Getenv("SUBSCRIPTION_CALLBACK_URL"),
		SuccessURL:  os.Getenv("SUBSCRIPTION_REDIRECT_URL"),
	}
	jsonReq, err := json.Marshal(invoiceReq)
	if err != nil {
		return err
	}
	utils.Logger.LogInfo().RawJSON("request", []byte(jsonReq)).Msg("json request data")

	// Create signature and add it to request headers
	base64Req := base64.StdEncoding.EncodeToString([]byte(jsonReq))

	sign := md5.Sum([]byte(base64Req + os.Getenv("GATEWAY_API_KEY")))

	utils.Logger.LogInfo().Str("sign", fmt.Sprintf("%x", sign)).Msg("Signature generated")

	req, err := http.NewRequest("POST", "https://api.cryptomus.com/v1/payment", strings.NewReader(string(jsonReq)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sign", fmt.Sprintf("%x", sign))
	req.Header.Set("merchant", os.Getenv("GATEWAY_MERCHANT"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	utils.Logger.LogInfo().Str("status", resp.Status).Msg("Invoice request sent")

	if resp.StatusCode != http.StatusOK {
		// Read body and return it
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return c.JSONBlob(resp.StatusCode, body)
	}
	// Parse response
	respStruct := struct {
		State  int `json:"state"`
		Result struct {
			Url     string `json:"url"`
			OrderID string `json:"order_id"`
			Uuid    string `json:"uuid"`
		} `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return err
	}

	utils.Logger.LogInfo().Fields(map[string]interface{}{
		"order_id": respStruct.Result.OrderID,
		"uuid":     respStruct.Result.Uuid,
	}).Msg("Invoice created")

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Invoice created",
		"invoice": map[string]any{
			"url":      respStruct.Result.Url,
			"order_id": respStruct.Result.OrderID,
			"uuid":     respStruct.Result.Uuid,
		},
	})
}

// ConfirmOrder is a webhook endpoint for cryptomus gateway
// It is called when payment is confirmed
// It checks if signature and order_id is valid and updates user subscription
func (contr *Controller) ConfirmOrder(c echo.Context) error {
	confirmReq := requests.ConfirmRequest{}
	if err := c.Bind(&confirmReq); err != nil {
		return err
	}
	// Verify signature
	jsonData, err := json.Marshal(confirmReq)
	if err != nil {
		return err
	}
	// Manually escape forward slashes
	escapedData := strings.ReplaceAll(string(jsonData), "/", "\\/")
	// Encode to base64
	var base64Req []byte
	base64.StdEncoding.Encode(base64Req, []byte(escapedData))
	// Compare hash
	hash := md5.Sum(append(base64Req, []byte(os.Getenv("GATEWAY_API_KEY"))...))
	if fmt.Sprintf("%x", hash) != confirmReq.Signature {
		utils.Logger.LogError().Fields(map[string]interface{}{
			"order_id": confirmReq.OrderID,
			"uuid":     confirmReq.Uuid,
			"status":   confirmReq.Status,
		}).Msg("Invalid signature")

		return fmt.Errorf("Invalid signature")
	}
	//Check status
	if confirmReq.Status != "paid" {
		utils.Logger.LogError().Fields(map[string]interface{}{
			"order_id": confirmReq.OrderID,
			"uuid":     confirmReq.Uuid,
			"status":   confirmReq.Status,
		}).Msg("Payment not confirmed")
		return fmt.Errorf("Payment not confirmed")
	}

	// Get user id from redis
	ctx := rediscl.RDB.Ctx
	userID, err := rediscl.RDB.Client.Get(ctx, "order_id:"+confirmReq.OrderID).Result()
	if err != nil {
		return err
	}
	// Update user subscription
	utils.Logger.LogInfo().Fields(map[string]interface{}{
		"order_id": confirmReq.OrderID,
		"user_id":  userID,
		"uuid":     confirmReq.Uuid,
		"status":   confirmReq.Status,
	}).Msg("Payment confirmed")
	// TODO: implement subscription update
	return c.JSON(http.StatusOK, map[string]any{
		"message":  "Payment confirmed",
		"order_id": confirmReq.OrderID,
		"uuid":     confirmReq.Uuid,
	})
}
