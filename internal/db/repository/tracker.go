package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"p2pbot/internal/db/models"
)

type TrackerRepository struct {
	db *sqlx.DB
}

func NewTrackerRepository(db *sqlx.DB) *TrackerRepository {
	return &TrackerRepository{db}
}

func (repo *TrackerRepository) UpdateWaitingFlag(id int64, flag bool) error {
	query := `UPDATE trackers SET waiting_adv = $1 WHERE id = $2`
	_, err := repo.db.Exec(query, flag, id)
	return err
}

func (repo *TrackerRepository) Save(tracker *models.Tracker) error {
	if tracker == nil {
		return fmt.Errorf("tracker is nil")
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	if tracker.ID == 0 {
		query := `INSERT INTO trackers (user_id, exchange, currency, side, waiting_adv)
			VALUES ($1, $2, $3, $4, false)
			ON CONFLICT(user_id, exchange, currency, side) DO UPDATE
			SET user_id = $1, exchange = $2, currency = $3, side = $4, waiting_adv = false
			RETURNING id`
		err := tx.QueryRow(query, tracker.ChatID, tracker.Exchange, tracker.Currency, tracker.Side).Scan(&tracker.ID)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error creating new tracker :", err)
		}
	} else {
		query := `UPDATE trackers SET exchange = $1, currency = $2, side = $3, waiting_adv = false
			WHERE id = $4`
		_, err = tx.Exec(query, tracker.Exchange, tracker.Currency, tracker.Side, tracker.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Remove old payment methods
		_, err = tx.Exec(`DELETE FROM methods WHERE tracker_id = $1`, tracker.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert new payment methods
	query := `INSERT INTO methods (tracker_id, payment_method)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING`
	for _, method := range tracker.Payment {
		_, err = tx.Exec(query, tracker.ID, method)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (repo *TrackerRepository) GetAllTrackers() ([]*models.UserTracker, error) {
	var trackers []*models.UserTracker
	query := `SELECT t.id, t.exchange, t.currency, t.side, t.waiting_adv, u.chat_id, u.binance_name, u.bybit_name
		FROM trackers t JOIN public.users u on t.user_id = u.chat_id`
	err := repo.db.Select(&trackers, query)
	if err != nil {
		return nil, err
	}

	for _, tracker := range trackers {
		rows, err := repo.db.Query("SELECT payment_method FROM methods WHERE tracker_id = $1", tracker.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting payment methods: %s", err)
		}

		for rows.Next() {
			var paymentMethod string
			err = rows.Scan(&paymentMethod)
			if err != nil {
				return nil, fmt.Errorf("error getting payment methods: %s", err)
			}

			tracker.Payment = append(tracker.Payment, paymentMethod)
		}
	}
	return trackers, nil
}
