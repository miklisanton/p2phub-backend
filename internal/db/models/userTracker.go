package models

type UserTracker struct {
    ID          int64    `db:"tracker_id" json:"id"`
	Exchange    string   `db:"exchange" json:"exchange"`
	Currency    string   `db:"currency" json:"currency"`
	Side        string   `db:"side" json:"side"`
	Waiting     bool     `db:"waiting_adv" json:"waiting"`
    Outbided bool        `db:"outbided" json:"outbided"`
	Payment     []PaymentMethod `db:"-" json:"payment_methods"`
    UserID      int       `db:"user_id" json:"-"`
	ChatID      *int64   `db:"chat_id" json:"tg_chat_id"`
	Username string      `db:"username" json:"username"`
}
