package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"p2pbot/internal/config"
	"strconv"
	"time"
)

type BybitExchange struct {
	adsEndpoint string
	name        string
	maxRetries  int
	retryDelay  time.Duration
}

type BybitPayload struct {
	TokenID    string   `json:"tokenId"`
	CurrencyID string   `json:"currencyId"`
	Side       string   `json:"side"`
	Payment    []string `json:"payment"`
	Page       string   `json:"page"`
}

type BybitAdsResponse struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	Result  Result `json:"result"`
}

type Result struct {
	Items []Item `json:"items"`
}

type Item struct {
	NickName          string   `json:"nickName"`
	Price             string   `json:"price"`
	Quantity          string   `json:"quantity"`
	MinAmount         string   `json:"minAmount"`
	MaxAmount         string   `json:"maxAmount"`
	Payments          []string `json:"payments"`
	RecentOrderNum    int      `json:"recentOrderNum"`
	RecentExecuteRate int      `json:"recentExecuteRate"`
}

func NewBybitExcahnge(config *config.Config) *BybitExchange {
	return &BybitExchange{
		adsEndpoint: "https://api2.bybit.com/fiat/otc/item/online",
		name:        "Bybit",
		maxRetries:  config.Exchange.MaxRetries,
		retryDelay:  time.Second * time.Duration(config.Exchange.RetryDelay),
	}
}

func (ex BybitExchange) GetName() string {
	return ex.name
}

func (ex BybitExchange) GetBestAdv(currency, side string, paymentMethods []string) (P2PItemI, error) {
	if side == "SELL" {
		side = "1"
	} else if side == "BUY" {
		side = "0"
	} else {
		return nil, fmt.Errorf("Invalid Side")
	}

	payload := BybitPayload{
		TokenID:    "USDT",
		CurrencyID: currency,
		Side:       side,
		Payment:    paymentMethods,
	}

	jsonPayload, err := json.Marshal(payload)

	resp, err := http.Post(ex.adsEndpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("could not make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	bybitResponse := BybitAdsResponse{}
	if err := json.Unmarshal(body, &bybitResponse); err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}

	if bybitResponse.RetCode != 0 {
		return nil, fmt.Errorf("bybit error: %s", bybitResponse.RetMsg)
	}

	return bybitResponse.Result.Items[0], nil
}

func (i Item) GetPrice() float64 {
	price, _ := strconv.ParseFloat(i.Price, 64)
	return price
}

func (i Item) GetQuantity() (quantity, minAmount, maxAmount float64) {
	quantity, _ = strconv.ParseFloat(i.Quantity, 64)
	minAmount, _ = strconv.ParseFloat(i.MinAmount, 64)
	maxAmount, _ = strconv.ParseFloat(i.MaxAmount, 64)
	return
}

func (i Item) GetName() string {
	return i.NickName
}

func (ex BybitExchange) requestData(page int, currency, side string) (*BybitAdsResponse, error) {
	if side == "SELL" {
		side = "1"
	} else if side == "BUY" {
		side = "0"
	} else {
		return nil, fmt.Errorf("invalid Side %s", side)
	}

	payload := BybitPayload{
		TokenID:    "USDT",
		CurrencyID: currency,
		Side:       side,
		Payment:    nil,
		Page:       fmt.Sprintf("%d", page),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("could marshal json: %w", err)
	}

	var body []byte

	for attempt := 1; attempt <= ex.maxRetries; attempt++ {
		resp, err := http.Post(ex.adsEndpoint, "application/json", bytes.NewBuffer(jsonPayload))
		if err == nil {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("bad status: %s", resp.Status)
				log.Println("retrying...")
				time.Sleep(ex.retryDelay)
				continue
			}

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("could not read response body: %v", err)
				log.Println("retrying...")
				time.Sleep(ex.retryDelay)
				continue
			}
			break
		}
		// sleep before retry
		if attempt < ex.maxRetries {
			time.Sleep(ex.retryDelay)
			log.Printf("could not connect to bybit exchange: %v, retrying...", err)
		} else {
			return nil, fmt.Errorf("could not connect to bybit exchange: %v, after %d attempts", err, ex.maxRetries)
		}
	}

	bybitResponse := BybitAdsResponse{}
	if err := json.Unmarshal(body, &bybitResponse); err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}

	if bybitResponse.RetCode != 0 {
		return nil, fmt.Errorf("bybit error: %s", bybitResponse.RetMsg)
	}
	return &bybitResponse, nil
}

func (ex BybitExchange) GetAdsByName(currency, side, username string) ([]P2PItemI, error) {
	out := make([]P2PItemI, 0)
	i := 1
	for {
		resp, err := ex.requestData(i, currency, side)
		if err != nil {
			return nil, fmt.Errorf("could not find advertisement with username %s", username)
		}

		if len(resp.Result.Items) == 0 {
			if len(out) == 0 {
				return nil, fmt.Errorf("could not find advertisement with username %s", username)
			} else {
				return out, nil
			}
		}

		for _, item := range resp.Result.Items {
			if item.GetName() == username {
				out = append(out, item)
			}
		}
		i++
	}
}

func (i Item) GetPaymentMethods() []string {
	return i.Payments
}
