package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // _ because we juts want to call this packages init function dont want to use its compnents
)

func Connect(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// setting for proper connection pool to prevent more db connections and to reuse database connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// fail fast
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// if db connection error return error like conection issue timeout erroes will catch here and return
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	return db, nil
}
