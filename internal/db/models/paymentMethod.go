package models

type PaymentMethod struct {
	Outbided   bool   `db:"outbidded"`
	tracker_id int64  `db:"tracker_id"`
	Id         string `db:"payment_method"`
	Name       string `db:"payment_name"`
}
