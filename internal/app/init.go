package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"p2pbot/internal/config"
	"p2pbot/internal/db/drivers"
	"p2pbot/internal/utils"
	"time"
)

func Init() (*sqlx.DB, *config.Config, error) {

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123Z}).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Caller().
		Logger().
		Hook(utils.GoroutineHook{})
	log.Info().Msg("Logger initialized")
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
