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

func (repo *UserRepository) Save(user *models.User) (int, error) {
    if user == nil {
        return 0, fmt.Errorf("user is nil")
    }



    if user.ID == 0 {
        query := `
            INSERT INTO users(email, chat_id) 
            VALUES ($1,$2)
            RETURNING id`

        row := repo.db.QueryRow(query,
                                    user.Email,
                                    user.ChatID)

        if row.Err() != nil {
            return 0, row.Err()
        }

        err := row.Scan(&user.ID)
        if err != nil {
            return 0, err
        }

        return user.ID, nil
    } else {
        query := `UPDATE users SET email = $1,
                    chat_id = $2
                    WHERE id = $3
                    RETURNING id`

        row := repo.db.QueryRow(query,
                                    user.Email,
                                    user.ChatID,
                                    user.ID)

        if row.Err() != nil {
            return 0, row.Err()
        }

        err := row.Scan(&user.ID)
        if err != nil {
            return 0, err
        }

        return user.ID, nil
    }

}

func (repo *UserRepository) Update(user *models.User) error {
    query := `UPDATE users SET email = $1,
                chat_id = $2
                WHERE id = $3`

    _, err := repo.db.Exec(query, user.Email,
                                    user.ChatID,
                                    user.ID)
    return err
}

func (repo *UserRepository) GetByChatID(chatID int64) (*models.User, error) {
    user := &models.User{}

    query := `SELECT * FROM users WHERE chat_id = $1`
    err := repo.db.Get(user, query, chatID)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (repo *UserRepository) GetByID(ID int) (*models.User, error) {
    user := &models.User{}

    query := `SELECT * FROM users WHERE id = $1`
    err := repo.db.Get(user, query, ID)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
    user := &models.User{}

    query := `SELECT * FROM users WHERE email = $1`
    err := repo.db.Get(user, query, email)
    if err != nil {
        return nil, err
    }

    return user, nil
}
