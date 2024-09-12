package models

type UserTracker struct {
    ID          int64    `db:"id" json:"id"`
	Exchange    string   `db:"exchange" json:"exchange"`
	Currency    string   `db:"currency" json:"currency"`
	Side        string   `db:"side" json:"side"`
	Waiting     bool     `db:"waiting_adv" json:"waiting"`
    Outbided bool        `db:"outbided" json:"outbided"`
	Payment     []string `db:"-" json:"payment_methods"`
	ChatID      *int64   `db:"chat_id" json:"tg_chat_id"`
	Username string      `db:"username" json:"username"`
}
