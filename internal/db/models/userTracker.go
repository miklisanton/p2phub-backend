package models

type UserTracker struct {
	ID          int64    `db:"id"`
	Exchange    string   `db:"exchange"`
	Currency    string   `db:"currency"`
	Side        string   `db:"side"`
	Waiting     bool     `db:"waiting_adv"`
	Payment     []string `db:"-"`
	ChatID      int64    `db:"chat_id"`
	BybitName   string   `db:"bybit_name"`
	BinanceName string   `db:"binance_name"`
}
