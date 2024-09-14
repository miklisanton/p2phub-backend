package models

type PaymentMethod struct {
    Outbided bool `db:"outbided"`
    tracker_id int64 `db:"tracker_id"`
    Name string `db:"name"`
    Notified bool `db:"notified"`
}
