package main

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatal("Failed to close database", zap.Error(err))
		}
	}(db)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal("Failed to configure database", zap.Error(err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		logger.Fatal("Failed to configure migration", zap.Error(err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal("Failed to run migration", zap.Error(err))
	}
}
