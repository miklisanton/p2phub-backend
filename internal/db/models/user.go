package models

type User struct {
	ChatID      int64  `db:"chat_id"`
	BinanceName string `db:"binance_name"`
	BybitName   string `db:"bybit_name"`
}
