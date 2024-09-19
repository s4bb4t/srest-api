package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

func SetupDataBase(dbStr string) (*Storage, error) {
	const op = "database.postgres.New"

	// Вывести текущую рабочую директорию
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("%s: unable to get current directory: %v", op, err)
	}
	fmt.Println("Current working directory:", dir)

	db, err := sql.Open("postgres", dbStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	migrationsDir := "../../internal/database/migrations"
	fmt.Println("Migrations directory:", migrationsDir) // Проверить путь

	if err := runMigrations(db, migrationsDir); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Storage{db: db}, nil
}

func runMigrations(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error setting postgres dialect: %v", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("error running migrations: %v", err)
	}

	return nil
}
