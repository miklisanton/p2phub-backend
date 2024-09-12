package models

type Tracker struct {
	ID       int64    `db:"id"`
	UserID   int    `db:"user_id"`
	Exchange string   `db:"exchange"`
	Currency string   `db:"currency"`
	Side     string   `db:"side"`
    Username string   `db:"username"`
	Waiting  bool     `db:"waiting_adv"`
    Outbided bool     `db:"outbided"`
	Payment  []string `db:"-"`
}
