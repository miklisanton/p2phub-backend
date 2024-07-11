package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"p2pbot/internal/db/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (repo *UserRepository) Save(user *models.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	query := `
		INSERT INTO users(chat_id, binance_name, bybit_name) 
		VALUES ($1, $2, $3)
		ON CONFLICT(chat_id) DO UPDATE SET 
		binance_name = coalesce(NULLIF(excluded.binance_name, ''), users.binance_name),
		bybit_name = coalesce(NULLIF(excluded.bybit_name, ''), users.bybit_name)`

	_, err := repo.db.Exec(query, user.ChatID, user.BinanceName, user.BybitName)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) GetByID(chatID int64) (*models.User, error) {
	user := &models.User{}

	query := `SELECT * FROM users WHERE chat_id = $1`
	err := repo.db.Get(user, query, chatID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
