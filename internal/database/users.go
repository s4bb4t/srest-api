package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
	"github.com/sabbatD/srest-api/internal/password"
)

func (s *Storage) Add(u u.User) (int, error) {
	const op = "database.postgres.Add"

	pwd, err := password.HashPassword(u.Password)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	var id int
	err = s.db.QueryRow(`
		INSERT INTO public.users (login, username, email, password, phone_number)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, u.Login, u.Username, u.Email, string(pwd), u.PhoneNumber).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" { // Код ошибки 23505 означает нарушение уникальности
			return 0, fmt.Errorf("%s: user already exists", op)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	_, err = s.db.Exec(`insert into public.roles (user_id, role) values ($1, ARRAY['USER'])`, id)

	if err != nil {
		return 0, fmt.Errorf("%s: INSERT INTO public.roles (user_id, role)\n\tvalues ($1, $2): %v", op, err)
	}

	return id, nil
}

func (s *Storage) Auth(u u.AuthData) (user u.TableUser, err error) {
	const op = "database.postgres.Auth"

	stmt, err := s.db.Prepare(`SELECT password FROM public.users WHERE login = $1`)
	if err != nil {
		return user, fmt.Errorf("%s.s.db.Prepare(`SELECT password FROM public.users WHERE login = $1`): %v", op, err)
	}
	defer stmt.Close()

	var pwd string

	if err = stmt.QueryRow(u.Login).Scan(&pwd); err != nil {
		return user, fmt.Errorf("%s.stmt.QueryRow(u.Login): %v", op, err)
	}

	if err := password.CheckPassword([]byte(pwd), u.Password); err != nil {
		return user, nil
	}

	stmt, err = s.db.Prepare(`SELECT id, username, email, date, is_blocked, is_admin FROM public.users WHERE login = $1`)
	if err != nil {
		return user, fmt.Errorf("%s.s.db.Prepare(`SELECT id, username, email, date, is_blocked, is_admin FROM public.users WHERE login = $1`): %v", op, err)
	}

	var isAdmin bool
	err = stmt.QueryRow(u.Login).Scan(&user.ID, &user.Username, &user.Email, &user.Date, &user.IsBlocked, &isAdmin)
	if err != nil {
		return user, fmt.Errorf("%s.stmt.QueryRow(u.Login).Scan(user): %v", op, err)
	}
	if user.IsBlocked {
		user.Roles = append(user.Roles, "USER")
	}

	err = s.db.QueryRow(`select role from public.roles where user_id = $1`, user.ID).Scan(pq.Array(&user.Roles))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if user.ID != 0 {
			_, err = s.db.Exec(`insert into public.roles (user_id, role) values ($1, ARRAY['USER'])`, user.ID)

			user.Roles = append(user.Roles, "USER")

			if err != nil {
				return user, fmt.Errorf("%s: %v", op, err)
			}

			return user, nil
		}
		return user, fmt.Errorf("%s/select role from roles where user_id = $1: %v", op, err)
	}

	return user, nil
}

func (s *Storage) UpdateRoles(id int, roles []string) (int64, error) {
	const op = "database.postgres.UpdateRoles"

	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM public.roles WHERE user_id = $1)`, id).Scan(&exists)
	if err != nil {
		return -1, fmt.Errorf("%s: %w while checking existence with user_id: %d", op, err, id)
	}

	if exists {
		res, err := s.db.Exec(`UPDATE public.roles SET role = $2 WHERE user_id = $1`, id, pq.Array(roles))
		if err != nil {
			return -1, fmt.Errorf("%s: %w while updating roles for user_id: %d", op, err, id)
		}

		n, err := res.RowsAffected()
		if err != nil {
			return -1, fmt.Errorf("%s: %w while fetching rows affected for user_id: %d", op, err, id)
		}
		return n, nil
	}

	res, err := s.db.Exec(`INSERT INTO public.roles (user_id, role) VALUES ($1, $2)`, id, pq.Array(roles))
	if err != nil {
		return -1, fmt.Errorf("%s: %w while inserting roles for user_id: %d", op, err, id)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("%s: %w while fetching rows affected for user_id: %d", op, err, id)
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
	defer stmt.Close()

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

func (s *Storage) All(q u.GetAllQuery) (result u.MetaResponse, E error) {
	const op = "database.postgres.GetAllUsers"

	qParams := []any{q.SearchTerm, q.Limit, q.Offset}
	mParams := []any{q.SearchTerm}

	query := `
		SELECT 
			id, 
			username, 
			email, 
			date, 
			is_blocked, 
			is_admin, 
			phone_number
		FROM public.users
		WHERE ($1 = '' OR username ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
		`
	metaQuery := `
		SELECT COUNT(*)
		FROM public.users
		WHERE ($1 = '' OR username ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
		`

	if q.IsBlocked != nil {
		query += ` AND is_blocked = $4 `
		metaQuery += ` AND is_blocked = $2 `

		qParams = append(qParams, *q.IsBlocked)
		mParams = append(mParams, *q.IsBlocked)
	}

	query += ` ORDER BY ` + q.SortBy + ` ` + q.SortOrder + ` LIMIT $2 OFFSET $3;`

	rows, err := s.db.Query(query, qParams...)
	if err != nil {
		return result, fmt.Errorf("%s: users req %v", op, err)
	}

	err = s.db.QueryRow(metaQuery, mParams...).Scan(&result.Meta.TotalAmount)
	if err != nil {
		return result, fmt.Errorf("%s: meta req %v", op, err)
	}

	defer rows.Close()

	var user u.TableUser
	var users []u.TableUser
	var isAdmin bool
	for rows.Next() {
		user.Roles = []string{"USER"}

		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Date, &user.IsBlocked, &isAdmin, &user.PhoneNumber); err != nil {
			return result, fmt.Errorf("%s: user scan %v", op, err)
		}

		if err := s.db.QueryRow(`select role from public.roles where user_id = $1`, user.ID).Scan(pq.Array(&user.Roles)); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("%s: roles scan %v", op, err)
		}

		users = append(users, user)
	}

	result.Meta.SortBy, result.Meta.SortOrder = q.SortBy, q.SortOrder
	result.Data = users

	return result, nil
}

func (s *Storage) Get(id int) (u.TableUser, error) {
	const op = "database.postgres.GetUser"

	var user u.TableUser
	var isAdmin bool

	err := s.db.QueryRow(`SELECT id, username, email, date, is_blocked, is_admin, phone_number FROM public.users WHERE id = $1`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Date, &user.IsBlocked, &isAdmin, &user.PhoneNumber)
	if err != nil {
		return u.TableUser{}, fmt.Errorf("%s: %v", op, err)
	}

	err = s.db.QueryRow(`select role from public.roles where user_id = $1;`, id).Scan(pq.Array(&user.Roles))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if user.ID != 0 {
			_, err = s.db.Exec(`insert into public.roles (user_id, role) values ($1, ARRAY['USER'])`, user.ID)

			user.Roles = append(user.Roles, "USER")

			if err != nil {
				return u.TableUser{}, fmt.Errorf("%s: %v", op, err)
			}

			return user, nil
		}
		return u.TableUser{}, fmt.Errorf("%s: %v", op, err)
	}

	return user, nil
}

func (s *Storage) UpdateField(field string, id int, val any) (int64, error) {
	const op = "database.postgres.UpdateUserField"

	switch field {
	case "admin":
		field = "is_admin"
	case "isadmin":
		field = "is_admin"
	case "IsAdmin":
		field = "is_admin"
	case "isAdmin":
		field = "is_admin"
	case "IsBlock":
		field = "is_admin"
	case "isblock":
		field = "is_admin"
	case "isBlock":
		field = "is_admin"
	case "block":
		field = "is_blocked"
	default:
		return -2, fmt.Errorf("%s: no such field: %v", op, field)
	}
	query := fmt.Sprintf(`UPDATE public.users SET %s = $1 WHERE id = $2`, field)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("%s: %v with parameters:%v, %v, %v", op, err, field, id, val)
	}
	defer stmt.Close()

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

func (s *Storage) UpdateUser(u u.PutUser, id int) (int64, error) {
	const op = "database.postgres.UpdateUser"

	tx, err := s.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	defer tx.Rollback()

	if u.Username != "" {
		_, err = tx.Exec(`UPDATE public.users SET username = $1 WHERE id = $2`, u.Username, id)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}
	}

	if u.Email != "" {
		var exists bool
		stmt, err := s.db.Prepare(`SELECT EXISTS (SELECT 1 FROM public.users WHERE email = $1)`)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}
		defer stmt.Close()

		if err = stmt.QueryRow(u.Email).Scan(&exists); err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}
		if exists {
			return -2, fmt.Errorf("%s: email already used", op)
		}

		_, err = tx.Exec(`UPDATE public.users SET email = $1 WHERE id = $2`, u.Email, id)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}
	}

	if u.PhoneNumber != "" {
		_, err = tx.Exec(`UPDATE public.users SET phone_number = $1 WHERE id = $2`, u.PhoneNumber, id)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	return 1, nil
}

func (s *Storage) SaveRefreshToken(token string, id int) error {
	const op = "database.postgres.SaveRefreshToken"

	stmt, err := s.db.Prepare(`
		INSERT INTO public.tokens (user_id, token, date) 
		VALUES ($1, $2, NOW() + INTERVAL '8 minutes') 
		ON CONFLICT (user_id) 
		DO UPDATE SET token = EXCLUDED.token, date = NOW() + INTERVAL '8 minutes'
	`)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, token)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (s *Storage) RefreshToken(token string) (string, int, error) {
	const op = "database.postgres.RefreshToken"

	stmt, err := s.db.Prepare(`SELECT user_id, token FROM public.tokens WHERE token = $1 and date > NOW()`)
	if err != nil {
		return "", 0, fmt.Errorf("%s: %v", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(token)
	if err != nil {
		return "", 0, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	if rows.Next() {
		var res string
		var id int
		if err := rows.Scan(&id, &res); err != nil {
			return "", 0, fmt.Errorf("%s: %v", op, err)
		}
		return token, id, nil
	}

	return "expired", 0, nil
}

func (s *Storage) ChangePassword(u u.Pwd, id int) (int64, error) {
	const op = "database.postgres.ChangePassword"

	var exists bool
	stmt, err := s.db.Prepare(`SELECT EXISTS (SELECT 1 FROM public.users WHERE id = $1)`)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}
	defer stmt.Close()

	if err = stmt.QueryRow(id).Scan(&exists); err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}
	if !exists {
		return -2, fmt.Errorf("%s: no such user", op)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	defer tx.Rollback()

	if u.Password != "" {
		pwd, err := password.HashPassword(u.Password)
		if err != nil {
			return 0, fmt.Errorf("%s: %v", op, err)
		}

		_, err = tx.Exec(`UPDATE public.users SET password = $1 WHERE id = $2`, pwd, id)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	return 1, nil
}

func (s *Storage) Logout(id int) error {
	const op = "database.postgres.Logout"

	res, err := s.db.Exec(`UPDATE public.users SET version = COALESCE(version, 0) + 1 WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	} else if n != 1 {
		return fmt.Errorf("%s: %s", op, "no rows affected")
	}

	_, err = s.db.Exec(`
		DELETE FROM public.tokens
		WHERE user_id = $1`, id)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (s *Storage) UserVersion(id int) int {
	const op = "database.postgres.UserVersion"

	stmt, err := s.db.Prepare(`SELECT version FROM public.users WHERE id = $1`)
	if err != nil {
		return 0
	}
	defer stmt.Close()

	ver := 0

	if err := stmt.QueryRow(id).Scan(&ver); err != nil {
		return 0
	}

	return ver
}
