package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
	pwd "github.com/sabbatD/srest-api/internal/password"
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

func (s *Storage) Add(u u.User) (int64, error) {
	const op = "database.postgres.Add"

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
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" { // Код ошибки 23505 означает нарушение уникальности
			return 0, fmt.Errorf("%s: user already exists", op)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s, failed to get lastInsertID: %v", op, err)
	}

	return id, nil
}

func (s *Storage) Auth(u u.AuthData) (succsess, admin bool, err error) {
	const op = "database.postgres.Auth"

	stmt, err := s.db.Prepare(`
		SELECT EXISTS (
			SELECT 1 FROM public.users 
			WHERE username = $1 AND password = $2
		)
	`)
	if err != nil {
		return true, false, fmt.Errorf("%s: %v", op, err)
	}

	hashedPassword, err := pwd.HashPassword(u.Password)
	if err != nil {
		return true, false, fmt.Errorf("%s: %v", op, err)
	}

	var exists bool
	if err = stmt.QueryRow(u.Username, string(hashedPassword)).Scan(&exists); err != nil {
		return false, false, fmt.Errorf("%s: %v", op, err)
	}

	stmt, err = s.db.Prepare(`SELECT admin, block FROM public.users WHERE username = $1`)
	if err != nil {
		return true, false, fmt.Errorf("%s: %v", op, err)
	}

	var isAdmin, isBlocked bool
	err = stmt.QueryRow(u.Username).Scan(&isAdmin, &isBlocked)
	if err != nil {
		return true, false, fmt.Errorf("%s: %v", op, err)
	}
	if isBlocked {
		isAdmin = false
	}
	return exists, isAdmin, nil
}

func (s *Storage) UpdateField(field string, u u.Login, val any) (int64, error) {
	const op = "database.postgres.UpdateField"
	query := fmt.Sprintf(`
	UPDATE public.users 
	SET %s = $1
	WHERE username = $2`, field)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, u.Username, val)
	}

	res, err := stmt.Exec(val, u.Username)
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, u.Username, val)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, u.Username, val)
	}

	if n == 0 {
		return n, fmt.Errorf("%s: no users with usernname: %v", op, u.Username)
	}

	return n, nil
}

func (s *Storage) Remove(u u.Login) (int64, error) {
	const op = "database.postgres.Remove"

	stmt, err := s.db.Prepare(`
	DELETE FROM public.users 
		WHERE username = $1
	`)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.Exec(u.Username)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	if n == 0 {
		return n, fmt.Errorf("%s: no users with usernname: %v", op, u.Username)
	}

	return n, nil
}

func (s *Storage) GetAll() ([]u.TableUser, error) {
	const op = "database.postgres.GetAll"

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

func (s *Storage) Get(username string) (u.TableUser, error) {
	const op = "database.postgres.Get"

	rows, err := s.db.Query(`SELECT * FROM public.users WHERE username = $1`, username)
	if err != nil {
		return u.TableUser{}, fmt.Errorf("%s: %v", op, err)
	}

	var user u.TableUser

	if rows.Next() {
		if err := rows.Scan(&user.Username, &user.Password, &user.Email, &user.Date, &user.Blocked, &user.Admin); err != nil {
			return u.TableUser{}, fmt.Errorf("%s: %v", op, err)
		}
	} else {
		return u.TableUser{}, fmt.Errorf("%s: no such user", op)
	}

	return user, nil
}

func (s *Storage) UpdateUser(u u.User, username string) (int64, error) {
	const op = "database.postgres.UpdateUser"

	stmt, err := s.db.Prepare(`
		SELECT EXISTS (
			SELECT 1 FROM public.users 
			WHERE username = $1 AND password = $2
		)
	`)
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v, %v", op, err, u.Username, u.Password, u.Email, username)
	}

	var exists bool
	if err = stmt.QueryRow(u.Username, u.Email).Scan(&exists); err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v, %v", op, err, u.Username, u.Password, u.Email, username)
	}
	if exists {
		return 0, fmt.Errorf("%s: username or email already used", op)
	}

	stmt, err = s.db.Prepare(`
		UPDATE public.users 
			SET username = $1, password = $2, email = $3
			WHERE username = $4
	`)
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v, %v", op, err, u.Username, u.Password, u.Email, username)
	}

	res, err := stmt.Exec(u.Username, u.Password, u.Email, username)
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v, %v", op, err, u.Username, u.Password, u.Email, username)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v, %v", op, err, u.Username, u.Password, u.Email, username)
	}

	if n == 0 {
		return n, fmt.Errorf("%s: no users with usernname: %v", op, username)
	}

	return n, nil
}

// TODO: forgot password help

// func (s *Storage) GetPass(u u.Login) (string, error) {

// }
