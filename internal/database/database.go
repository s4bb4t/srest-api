package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

var DB *sql.DB

func SetupDataBase(dbStr, env string) (*Storage, error) {
	const op = "database.postgres.New"

	db, err := sql.Open("postgres", dbStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	if env == "local" {
		migrationsDir := "./internal/database/migrations"
		fmt.Println("Migrations directory:", migrationsDir) // Проверить путь

		if err := runMigrations(db, migrationsDir); err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
	}

	DB = db
	return &Storage{db: db}, nil
}

func runMigrations(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error setting postgres dialect: %v", err)
	}

	// Сброс миграций
	// if err := goose.Reset(db, migrationsDir); err != nil {
	// 	return fmt.Errorf("error resetting migrations: %v", err)
	// }

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("error running migrations: %v", err)
	}

	return nil
}

func UserVersion(id int) int {
	const op = "database.postgres.UserVersion"

	stmt, err := DB.Prepare(`SELECT version FROM public.users WHERE id = $1`)
	if err != nil {
		return 0
	}
	defer stmt.Close()

	var ver int

	if err := stmt.QueryRow(id).Scan(&ver); err != nil {
		fmt.Println(err.Error())
		return 0
	}
	fmt.Println("UserVer from access", id, ver)
	return ver
}
