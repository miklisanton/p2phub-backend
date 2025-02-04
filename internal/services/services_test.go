package services_test

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"p2pbot/internal/config"
	"p2pbot/internal/services"
	"path/filepath"
	"testing"
	"time"
)

var (
	bybit   *services.BybitExchange
	binance *services.BinanceExchange
)

func findProjectRoot() (string, error) {
	// Get the current working directory of the test
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse upwards to find the project root (where .env file is located)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); os.IsNotExist(err) {
			parent := filepath.Dir(dir)
			if parent == dir {
				// Reached the root of the filesystem, and .env wasn't found
				return "", os.ErrNotExist
			}
			dir = parent
		} else {
			return dir, nil
		}
	}
}

func TestMain(m *testing.M) {
	root, err := findProjectRoot()
	fmt.Println("Root: ", root)
	if err != nil {
		panic("Error finding project root: " + err.Error())
	}

	err = os.Chdir(root)
	if err != nil {
		panic(err)
	}

	cfg, err := config.NewConfig("/Users/antonmiklis/GolandProjects/p2pbot/config.yaml")
	if err != nil {
		panic(err)
	}
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123Z}).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	bybit = services.NewBybitExcahnge(cfg)
	code := m.Run()

	//DB.MustExec("DELETE FROM users");

	os.Exit(code)
}
func TestBybitGetCachedMethods(t *testing.T) {
	currencies, err := bybit.GetCachedPaymentMethods("CZK")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Log("result: ", currencies)
}

//func TestBinanceFetchMethods(t *testing.T) {
//    currencies, err := binance.FetchCurrencies()
//    if err != nil {
//        t.Errorf("Error: %v", err)
//    }
//
//    res, err := binance.FetchPaymentMethods(currencies)
//    if err != nil {
//        t.Errorf("Error: %v", err)
//    }
//
//    t.Log("result: ", res)
//}

func TestBinanceGetCachedMethods(t *testing.T) {
	currencies, err := binance.GetCachedPaymentMethods("")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Log("result: ", currencies)
}

func TestBinanceGetCachedCurrencies(t *testing.T) {
	currencies, err := binance.GetCachedCurrencies()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Log("result: ", currencies)
}

func TestBybitGetCachedCurrencies(t *testing.T) {
	currencies, err := bybit.GetCachedCurrencies()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Log("result: ", currencies)
}

func TestBybitFetchAds(t *testing.T) {
	t.Log("Fetching ads")
	ads, err := bybit.GetAds("EUR", "BUY")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Log("result: ", ads)
}
