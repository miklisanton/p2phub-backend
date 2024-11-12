package repository

import (
	"database/sql"
	"fmt"
	"p2pbot/internal/db/models"

	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (repo *SubscriptionRepository) Save(subscription *models.Subscription) error {
	if subscription == nil {
		return fmt.Errorf("subscription is nil")
	}

	if subscription.Id != 0 {
		query := `UPDATE subscription SET user_id = $1, created_at = $2, valid_until = $3 WHERE id = $4`
		_, err := repo.db.Exec(query, subscription.User_id, subscription.Created_at, subscription.ValidUntil, subscription.Id)
		if err != nil {
			return fmt.Errorf("error updating subscription : %v", err)
		}
		return nil
	} else {
		query := `INSERT INTO subscription (user_id, created_at, valid_until) VALUES ($1, $2, $3) RETURNING id`
		err := repo.db.QueryRow(query, subscription.User_id, subscription.Created_at, subscription.ValidUntil).Scan(&subscription.Id)
		if err != nil {
			return fmt.Errorf("error creating new subscription : %v", err)
		}
	}

	return nil
}

func (repo *SubscriptionRepository) GetByID(id int) (*models.Subscription, error) {
	query := `SELECT * FROM subscription WHERE id = $1`
	var subscription models.Subscription
	err := repo.db.Get(&subscription, query, id)
	if err != nil {
		return nil, fmt.Errorf("error getting subscription by id : %v", err)
	}

	return &subscription, nil
}

func (repo *SubscriptionRepository) GetByUserID(id int) (*models.Subscription, error) {
	query := `SELECT * FROM subscription WHERE user_id = $1`
	var subscription models.Subscription
	err := repo.db.Get(&subscription, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting subscription by user id : %v", err)
	}
	return &subscription, nil
}
