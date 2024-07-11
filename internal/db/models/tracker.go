package models

type Tracker struct {
	ID       int64    `db:"id"`
	ChatID   int64    `db:"user_id"`
	Exchange string   `db:"exchange"`
	Currency string   `db:"currency"`
	Side     string   `db:"side"`
	Waiting  bool     `db:"waiting_adv"`
	Payment  []string `db:"-"`
}
