package models

type PaymentMethod struct {
    Outbided bool `db:"outbided"`
    tracker_id int64 `db:"tracker_id"`
    Id string `db:"name"`
    Name string `db:"payment_name"`
    Notified bool `db:"notified"`
}
