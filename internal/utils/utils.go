package utils

import (
	"fmt"
    "p2pbot/internal/db/models"
	"p2pbot/internal/services"
	"reflect"
    "encoding/json"
    "golang.org/x/crypto/bcrypt"
)

type Notification struct {
    ChatID    int64                 `json:"chat_id"`
	Data      services.P2PItemI     `json:"top_order"`
	Exchange  string                `json:"exchange"`
	Side string                     `json:"side"` 
	Currency  string                `json:"currency"`
}

func (n *Notification) UnmarshalJSON(data []byte) error {
    // Define a temporary structure for the concrete type
    type Alias Notification
    aux := &struct {
        Data json.RawMessage `json:"top_order"`
        *Alias
    }{
        Alias: (*Alias)(n),
    }

    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    if aux.Exchange == "binance" {
        var item services.DataItem
        if err := json.Unmarshal(aux.Data, &item); err != nil {
            return err
        }
        n.Data = item
    } else if aux.Exchange == "bybit"{
        var item services.Item
        if err := json.Unmarshal(aux.Data, &item); err != nil {
            return err
        }
        n.Data = item
    } else {
        return fmt.Errorf("no unmarshal logic for %s", aux.Exchange)
    }

    return nil
}

func GetField(obj interface{}, name string) (interface{}, error) {
	v := reflect.ValueOf(obj).Elem()

	field := v.FieldByName(name)
	if !field.IsValid() {
		return nil, fmt.Errorf("no such field: %s in obj", name)
	}

	return field.Interface(), nil
}

func SetField(obj interface{}, name string, value interface{}) error {
	v := reflect.ValueOf(obj).Elem()

	field := v.FieldByName(name)
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	field.Set(reflect.ValueOf(value))
	return nil
}

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

// Returns true if one of trackers payment methods is in ads payment methods
func ComparePaymentMethods(adMethods []string, trackerMethods []*models.PaymentMethod) bool {
    for _, trackerMethod := range trackerMethods {
        if Contains(adMethods, trackerMethod.Id) {
            return true
        } 
    }
    return false
}

func AllOutbidded(pMethods []*models.PaymentMethod) bool {
    for _, pMethod := range pMethods {
        if !pMethod.Outbided {
            return false
        }
    }
    return true
}

func GetPMethodName(pMethods []services.PaymentMethod, id string) (string, error) {
    for _, pMethod := range pMethods {
        if pMethod.Id == id {
            return pMethod.Name, nil
        }
    }
    return "", fmt.Errorf("payment method not found")
}
