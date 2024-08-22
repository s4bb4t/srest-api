package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	pwd "github.com/sabbatD/srest-api/internal/password"
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

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS public.users (
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

	stmt, err := s.db.Prepare(`
		INSERT INTO public.users (
			username, email, password, date
		) VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	hashedPassword, err := pwd.HashPassword(u.Password)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.Exec(u.Username, u.Email, string(hashedPassword), time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s, failed to get lastInsertID: %v", op, err)
	}

	return id, nil
}

func (s *Storage) Auth(u u.AuthData) (bool, error) {
	const op = "database.postgres.Auth"

	stmt, err := s.db.Prepare(`
		SELECT EXISTS (
			SELECT 1 FROM public.users 
			WHERE username = $1 AND password = $2
		)
	`)
	if err != nil {
		return false, fmt.Errorf("%s: %v", op, err)
	}

	hashedPassword, err := pwd.HashPassword(u.Password)
	if err != nil {
		return false, fmt.Errorf("%s: %v", op, err)
	}

	var exists bool
	if err = stmt.QueryRow(u.Username, string(hashedPassword)).Scan(&exists); err != nil {
		return false, fmt.Errorf("%s: %v", op, err)
	}

	return exists, nil
}

func (s *Storage) UpdateField(field string, u u.Login, val any) error {
	const op = "database.postgres.UpdateField"

	stmt, err := s.db.Prepare(`
		UPDATE public.users 
			SET $1 = $2
			WHERE username = $3
	`)
	if err != nil {
		return fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, u.Username, val)
	}

	res, err := stmt.Exec(field, val, u.Username)
	if err != nil {
		return fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, u.Username, val)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, u.Username, val)
	}

	if n == 0 {
		return fmt.Errorf("%s: no users with usernname: %v", op, u.Username)
	}

	return nil
}

func (s *Storage) RemoveUser(u u.Login) error {
	const op = "database.postgres.RemoveUser"

	stmt, err := s.db.Prepare(`
	DELETE FROM public.users 
		WHERE username = $1
	`)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.Exec(u.Username)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if n == 0 {
		return fmt.Errorf("%s: no users with usernname: %v", op, u.Username)
	}

	return nil
}

func (s *Storage) GetAllUsers() ([]u.TableUser, error) {
	const op = "database.postgres.GetAllUsers"

	rows, err := s.db.Query(`SELECT * FROM public.users`)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	var result []u.TableUser
	var user u.TableUser

	for rows.Next() {
		if err := rows.Scan(&user.Username, &user.Password, &user.Email, &user.Date, &user.Blocked, &user.Admin); err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}

		result = append(result, user)
	}

	return result, nil
}

// TODO: forgot password help

// func (s *Storage) GetPass(u u.Login) (string, error) {

// }
