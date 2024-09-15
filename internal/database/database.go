package database

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func SetupDataBase(dbStr string) (*Storage, error) {
	const op = "database.postgres.New"

	db, err := sql.Open("postgres", dbStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	// uuid
	stmt, err := db.Prepare(` 
		CREATE TABLE IF NOT EXISTS public.users (
			id INTEGER PRIMARY KEY UNIQUE, 
			login TEXT UNIQUE,
			username TEXT,
			password TEXT,
			email TEXT UNIQUE,
			date TEXT,
			is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE,
			phone BIGINT
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS public.todos (
			id INTEGER PRIMARY KEY UNIQUE,
			title TEXT,
			created TEXT,
			is_done BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	// tz
	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS public.tokens (
			user_id INTEGER PRIMARY KEY UNIQUE,
			token TEXT,
			date TIMESTAMP DEFAULT NOW() + INTERVAL '1 day' 
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Storage{db: db}, nil
}
