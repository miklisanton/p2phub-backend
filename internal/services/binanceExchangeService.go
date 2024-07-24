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

type BinanceExchange struct {
	adsEndpoint string
	name        string
	maxRetries  int
	retryDelay  time.Duration
}

type BinancePayload struct {
	Fiat                      string   `json:"fiat"`
	Page                      int      `json:"page"`
	Rows                      int      `json:"rows"`
	TradeType                 string   `json:"tradeType"`
	Asset                     string   `json:"asset"`
	Countries                 []string `json:"countries"`
	ProMerchantAds            bool     `json:"proMerchantAds"`
	ShieldMerchantAds         bool     `json:"shieldMerchantAds"`
	FilterType                string   `json:"filterType"`
	Periods                   []string `json:"periods"`
	AdditionalKycVerifyFilter int      `json:"additionalKycVerifyFilter"`
	PublisherType             *string  `json:"publisherType"`
	PayTypes                  []string `json:"payTypes"`
	Classifies                []string `json:"classifies"`
}

type BinanceAdsResponse struct {
	Code    string     `json:"code"`
	Data    []DataItem `json:"data"`
	Total   int        `json:"total"`
	Success bool       `json:"success"`
}

type DataItem struct {
	Adv        Adv        `json:"adv"`
	Advertiser Advertiser `json:"advertiser"`
}

type Adv struct {
	Price                string        `json:"price"`
	TradableQuantity     string        `json:"tradableQuantity"`
	MaxSingleTransAmount string        `json:"maxSingleTransAmount"`
	MinSingleTransAmount string        `json:"minSingleTransAmount"`
	TradeMethods         []TradeMethod `json:"tradeMethods"`
	IsTradable           bool          `json:"isTradable"`
}

type TradeMethod struct {
	Identifier string `json:"identifier"`
}

type Advertiser struct {
	NickName        string  `json:"nickName"`
	MonthOrderCount int     `json:"monthOrderCount"`
	MonthFinishRate float64 `json:"monthFinishRate"`
	PositiveRate    float64 `json:"positiveRate"`
}

func NewBinanceExchange(config *config.Config) *BinanceExchange {
	return &BinanceExchange{
		adsEndpoint: "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search",
		name:        "Binance",
		maxRetries:  config.Exchange.MaxRetries,
		retryDelay:  time.Second * time.Duration(config.Exchange.RetryDelay),
	}
}

func (ex BinanceExchange) GetName() string {
	return ex.name
}

func (ex BinanceExchange) GetBestAdv(currency, side string, paymentMethods []string) (P2PItemI, error) {
	if side == "BUY" {
		side = "SELL"
	} else if side == "SELL" {
		side = "BUY"
	} else {
		return nil, fmt.Errorf("invalid side")
	}

	payload := BinancePayload{
		Fiat:                      currency,
		Page:                      1,
		Rows:                      10,
		TradeType:                 side,
		Asset:                     "USDT",
		Countries:                 []string{},
		ProMerchantAds:            false,
		ShieldMerchantAds:         false,
		FilterType:                "all",
		Periods:                   []string{},
		AdditionalKycVerifyFilter: 0,
		PublisherType:             nil,
		PayTypes:                  paymentMethods,
		Classifies:                []string{"mass", "profession"},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
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
			log.Printf("could not connect to binance exchange: %v, retrying...", err)
		} else {
			return nil, fmt.Errorf("could not connect to binance exchange: %v, after %d attempts", err, ex.maxRetries)
		}
	}

	binanceResponse := BinanceAdsResponse{}
	if err := json.Unmarshal(body, &binanceResponse); err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}

	if !binanceResponse.Success {
		return nil, fmt.Errorf("binance error: %s", binanceResponse.Code)
	}

	if len(binanceResponse.Data) == 0 {
		log.Println(payload)

		return nil, fmt.Errorf("binance response has no data")
	}
	return binanceResponse.Data[0], nil
}

func (i DataItem) GetName() string {
	return i.Advertiser.NickName
}

func (i DataItem) GetPrice() float64 {
	price, _ := strconv.ParseFloat(i.Adv.Price, 64)
	return price
}

func (i DataItem) GetQuantity() (quantity, minAmount, maxAmount float64) {
	quantity, _ = strconv.ParseFloat(i.Adv.TradableQuantity, 64)
	minAmount, _ = strconv.ParseFloat(i.Adv.MinSingleTransAmount, 64)
	maxAmount, _ = strconv.ParseFloat(i.Adv.MaxSingleTransAmount, 64)
	return
}

func (ex BinanceExchange) RequestData(page int, currency, side string) (*BinanceAdsResponse, error) {
	if side == "BUY" {
		side = "SELL"
	} else if side == "SELL" {
		side = "BUY"
	} else {
		return nil, fmt.Errorf("invalid side")
	}

	payload := BinancePayload{
		Fiat:                      currency,
		Page:                      page,
		Rows:                      10,
		TradeType:                 side,
		Asset:                     "USDT",
		Countries:                 []string{},
		ProMerchantAds:            false,
		ShieldMerchantAds:         false,
		FilterType:                "all",
		Periods:                   []string{},
		AdditionalKycVerifyFilter: 0,
		PublisherType:             nil,
		PayTypes:                  make([]string, 0),
		Classifies:                []string{"mass", "profession"},
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

	binanceResponse := BinanceAdsResponse{}
	if err := json.Unmarshal(body, &binanceResponse); err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}

	if !binanceResponse.Success {
		return nil, fmt.Errorf("binance error: %s", binanceResponse.Code)
	}

	return &binanceResponse, nil
}
func (ex BinanceExchange) GetAdsByName(currency, side, username string) ([]P2PItemI, error) {
	out := make([]P2PItemI, 0)
	i := 1
	for {
		response, err := ex.RequestData(i, currency, side)
		if err != nil {

			return nil, fmt.Errorf("could not find advertisement with username %s", username)
		}
		// All pages parsed, adv not found
		if len(response.Data) == 0 {
			if len(out) == 0 {
				return nil, fmt.Errorf("could not find advertisement with username %s", username)
			} else {
				return out, nil
			}
		}

		for _, item := range response.Data {
			if item.GetName() == username {
				out = append(out, item)
			}
		}
		i++
	}
}

func (i DataItem) GetPaymentMethods() []string {
	out := make([]string, 0)
	for _, method := range i.Adv.TradeMethods {
		out = append(out, method.Identifier)
	}
	return out
}
