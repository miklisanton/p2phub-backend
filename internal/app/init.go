package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"p2pbot/internal/config"
	"p2pbot/internal/db/drivers"
)

func Init() (*sqlx.DB, *config.Config, error) {

	cfg, err := config.NewConfig("config.yaml")
	if err != nil {
		return nil, nil, err
	}

	connectionURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host+":"+cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSL)
	DB, err := drivers.Connect(connectionURL)
	if err != nil {
		return nil, nil, err
	}

	return DB, cfg, nil
}
