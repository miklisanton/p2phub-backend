package services_test

import (
    "fmt"
    "os"
    "p2pbot/internal/app"
    "p2pbot/internal/rediscl"
    "p2pbot/internal/services"
    "path/filepath"
    "testing"
)

var (
    bybit *services.BybitExchange
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

    _, cfg, err := app.Init()
    if err != nil {
        panic(err)
    }
    bybit = services.NewBybitExcahnge(cfg)
    binance = services.NewBinanceExchange(cfg)
    rediscl.InitRedisClient(cfg.Redis.Host, cfg.Redis.Port)
    code := m.Run()

    //DB.MustExec("DELETE FROM users");
    
    os.Exit(code)
}

//func TestFetchAllPaymentList(t *testing.T) {
//    res, err := bybit.FetchAllPaymentList()
//    if err != nil {
//        t.Errorf("Error: %v", err)
//    }
//    t.Log("result: ", res)
//}

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
