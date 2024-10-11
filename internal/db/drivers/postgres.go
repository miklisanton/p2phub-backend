package drivers

import (
    "database/sql"
    // import go libpq driver package
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/pressly/goose/v3"
)

func Connect(connectionURL string) (*sqlx.DB, error) {
    db, err := sql.Open("postgres", connectionURL)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    if err = runMigrations(db); err != nil {
        return nil, err
    }

    sqlxDB := sqlx.NewDb(db, "postgres")
    return sqlxDB, nil
}

func runMigrations(db *sql.DB) error {
    goose.SetDialect("postgres")

    if err := goose.Up(db, "./internal/db/migrations"); err != nil {
        return err
    }
    return nil
}
