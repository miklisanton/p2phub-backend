package models

import (
	"time"
)

type Subscription struct {
	Id          int       `db:"id" json:"id"`
	User_id     int       `db:"user_id" json:"user_id"`
	Created_at  time.Time `db:"created_at" json:"created_at"`
	Valid_until time.Time `db:"valid_until" json:"valid_until"`
}
