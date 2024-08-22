package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	u "github.com/sabbatD/srest-api/internal/user"
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

	// stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (username TEXT UNIQUE password TEXT email TEXT UNIQUE date TEXT block BOOLEAN NOT NULL DEFAULT FALSE admin BOOLEAN NOT NULL DEFAULT FALSE)")
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			username TEXT UNIQUE,
			password TEXT,
			email TEXT UNIQUE,
			date TEXT,
			block BOOLEAN NOT NULL DEFAULT FALSE,
			admin BOOLEAN NOT NULL DEFAULT FALSE
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

func (s *Storage) AddNewUser(u u.User) (int64, error) {
	const op = "database.postgres.AddNewUser"

	stmt, err := s.db.Prepare("INSERT INTO public.users (username, email, password, date) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.Exec(u.Username, u.Email, u.Password, time.Now())
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s, failed to get lastInsertID: %v", op, err)
	}

	return id, nil
}
