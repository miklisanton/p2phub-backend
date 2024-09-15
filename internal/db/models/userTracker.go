package models

type UserTracker struct {
    ID          int64    `db:"tracker_id" json:"id"`
	Exchange    string   `db:"exchange" json:"exchange"`
	Currency    string   `db:"currency" json:"currency"`
	Side        string   `db:"side" json:"side"`
    Notify      bool     `db:"notify" json:"notify"`
	Payment     []PaymentMethod `db:"-" json:"payment_methods"`
    UserID      int       `db:"user_id" json:"-"`
	ChatID      *int64   `db:"chat_id" json:"tg_chat_id"`
	Username string      `db:"username" json:"username"`
}
