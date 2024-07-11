package utils

import (
	"fmt"
	"p2pbot/internal/services"
	"reflect"
)

type Notification struct {
	ChatID    int64
	Data      services.P2PItemI
	Exchange  string
	Direction string
	Currency  string
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
