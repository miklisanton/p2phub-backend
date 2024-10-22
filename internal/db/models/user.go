package models

type User struct {
	ID     int     `db:"id"`
	ChatID *int64  `db:"chat_id"`
	Email  *string `db:"email"`
}
