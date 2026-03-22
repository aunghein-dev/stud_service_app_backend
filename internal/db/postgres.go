package db

import (
	"database/sql"
	"fmt"
	"time"

	"student_service_app/backend/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgres(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}
