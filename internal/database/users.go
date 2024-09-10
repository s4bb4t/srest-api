package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
)

func (s *Storage) Add(u u.User) (int, error) {
	const op = "database.postgres.Add"

	stmt, err := s.db.Prepare(`
		INSERT INTO public.users (
			id, login, username, email, password, date
		) VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	var maxID sql.NullInt64

	err = s.db.QueryRow(`SELECT MAX(id) FROM public.users;`).Scan(&maxID)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	if !maxID.Valid {
		maxID.Int64 = 1
	} else {
		maxID.Int64 += 1
	}

	res, err := stmt.Exec(maxID.Int64, u.Login, u.Username, u.Email, u.Password, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" { // Код ошибки 23505 означает нарушение уникальности
			return 0, fmt.Errorf("%s: user already exists", op)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		return 0, fmt.Errorf("%s, failed to get RowsAffected: %v", op, err)
	}

	return int(maxID.Int64), nil
}

func (s *Storage) Auth(u u.AuthData) (user u.TableUser, err error) {
	const op = "database.postgres.Auth"

	stmt, err := s.db.Prepare(`
		SELECT EXISTS (
			SELECT 1 FROM public.users 
			WHERE login = $1 AND password = $2
		)
	`)
	if err != nil {
		return user, fmt.Errorf("%s: %v", op, err)
	}

	var exists bool
	if err = stmt.QueryRow(u.Login, u.Password).Scan(&exists); err != nil {
		return user, fmt.Errorf("%s: %v", op, err)
	}

	if exists {
		stmt, err = s.db.Prepare(`SELECT * FROM public.users WHERE login = $1`)
		if err != nil {
			return user, fmt.Errorf("%s: %v", op, err)
		}

		err = stmt.QueryRow(u.Login).Scan(&user.ID, &user.Username, &user.Login, &user.Password, &user.Email, &user.Date, &user.Block, &user.Admin)
		if err != nil {
			return user, fmt.Errorf("%s: %v", op, err)
		}
		if user.Block {
			user.Admin = false
		}
	}

	return user, nil
}

func (s *Storage) UpdateField(field string, id int, val any) (int64, error) {
	const op = "database.postgres.UpdateUserField"

	switch field {
	case "admin":
	case "block":
	default:
		return 0, fmt.Errorf("%s: no such field: %v", op, field)
	}
	query := fmt.Sprintf(`UPDATE public.users SET %s = $1 WHERE id = $2`, field)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, id, val)
	}

	res, err := stmt.Exec(val, id)
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, id, val)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, id, val)
	}

	if n == 0 {
		return n, fmt.Errorf("%s: no users with id: %v", op, id)
	}

	return n, nil
}

func (s *Storage) Remove(id int) (int64, error) {
	const op = "database.postgres.RemoveUser"

	stmt, err := s.db.Prepare(`
	DELETE FROM public.users 
		WHERE id = $1
	`)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	if n == 0 {
		return n, fmt.Errorf("%s: no users with id: %v", op, id)
	}

	return n, nil
}

func (s *Storage) GetAll(search, order string, blocked bool, limit, offset int) ([]u.TableUser, error) {
	const op = "database.postgres.GetAllUsers"

	query := ``
	if order == "none" {
		order = "id"
		query = `
			SELECT * FROM public.users
			WHERE ($1 = '' OR login ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%' OR username ILIKE '%' || $1 || '%') AND block = $2
			ORDER BY $3 ASC
			LIMIT $4 OFFSET $5
		`
	} else {
		query = `
			SELECT * FROM public.users
			WHERE (($1 = '' OR username ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%') AND block = $2)
			ORDER BY CASE WHEN $3 = 'asc' THEN email END ASC,
					CASE WHEN $3 = 'desc' THEN email END DESC
			LIMIT $4 OFFSET $5
		`
	}

	rows, err := s.db.Query(query, search, blocked, order, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	var result []u.TableUser
	var user u.TableUser

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Login, &user.Username, &user.Password, &user.Email, &user.Date, &user.Block, &user.Admin); err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}

		result = append(result, user)
	}

	return result, nil
}

func (s *Storage) Get(id int) (u.TableUser, error) {
	const op = "database.postgres.GetUser"

	rows, err := s.db.Query(`SELECT * FROM public.users WHERE id = $1`, id)
	if err != nil {
		return u.TableUser{}, fmt.Errorf("%s: %v", op, err)
	}

	var user u.TableUser

	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.Login, &user.Username, &user.Password, &user.Email, &user.Date, &user.Block, &user.Admin); err != nil {
			return u.TableUser{}, fmt.Errorf("%s: %v", op, err)
		}
	} else {
		return u.TableUser{}, fmt.Errorf("%s: no such user", op)
	}

	return user, nil
}

func (s *Storage) UpdateUser(u u.User, id int) (int64, error) {
	const op = "database.postgres.UpdateUser"

	stmt, err := s.db.Prepare(`
		SELECT EXISTS (
			SELECT 1 FROM public.users 
			WHERE login = $1 OR email = $2
		)
	`)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	var exists bool
	if err = stmt.QueryRow(u.Login, u.Email).Scan(&exists); err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}
	if exists {
		return 0, fmt.Errorf("%s: login or email already used", op)
	}

	stmt, err = s.db.Prepare(`
		UPDATE public.users 
			SET login = $1, username = $2, password = $3, email = $4
			WHERE id = $5
	`)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.Exec(u.Login, u.Username, u.Password, u.Email, id)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	if n == 0 {
		return n, fmt.Errorf("%s: no users with id: %v", op, id)
	}

	return n, nil
}

func (s *Storage) SaveRefreshToken(token string, id int) error {
	const op = "database.postgres.SaveRefreshToken"

	stmt, err := s.db.Prepare(`INSERT INTO public.tokens (user_id, token) VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	_, err = stmt.Exec(id, token)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
func (s *Storage) RefreshToken(token string) (string, int, error) {
	const op = "database.postgres.SaveRefreshToken"

	stmt, err := s.db.Prepare(`SELECT user_id, token FROM public.tokens WHERE token = $1 and date < Now()`)
	if err != nil {
		return "", 0, fmt.Errorf("%s: %v", op, err)
	}

	rows, err := stmt.Query(token)
	if err != nil {
		return "", 0, fmt.Errorf("%s: %v", op, err)
	}

	var res string
	var id int

	for !rows.Next() {
		if err := rows.Scan(&id, &token); err != nil {
			return "", 0, fmt.Errorf("%s: %v", op, err)
		}
	}

	if res == "" {
		res = "expired"
	}

	return token, id, nil
}
