package database

import (
	"database/sql"
	"fmt"
	"time"

	t "github.com/sabbatD/srest-api/internal/lib/todoConfig"
)

func (s *Storage) Create(t t.TodoRequest) (int64, error) {
	const op = "database.postgres.CreateTodo"

	stmt, err := s.db.Prepare(`
		INSERT INTO public.todos (
			id, title, created, is_done
		) VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	var maxID sql.NullInt64

	err = s.db.QueryRow(`SELECT MAX(id) FROM public.todos;`).Scan(&maxID)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	if !maxID.Valid {
		maxID.Int64 = 1
	} else {
		maxID.Int64 += 1
	}

	_, err = stmt.Exec(maxID.Int64, t.Title, time.Now().Format("2006-01-02 15:04:05"), *t.IsDone)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return maxID.Int64, nil
}

func (s *Storage) Update(id int, t t.TodoRequest) (int64, error) {
	const op = "database.postgres.UpdateTodo"

	var stmt *sql.Stmt
	var err error
	var res sql.Result

	if t.Title == "" {
		stmt, err = s.db.Prepare(`UPDATE public.todos SET is_done = $1 WHERE id = $2`)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}

		res, err = stmt.Exec(*t.IsDone, id)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}

	}

	if t.IsDone == nil {
		stmt, err = s.db.Prepare(`UPDATE public.todos SET title = $1 WHERE id = $2`)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}
		res, err = stmt.Exec(t.Title, id)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}

	}

	if t.Title != "" && t.IsDone != nil {
		stmt, err = s.db.Prepare(`UPDATE public.todos SET title = $1, is_done = $2 WHERE id = $3`)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}
		res, err = stmt.Exec(t.Title, t.IsDone, id)
		if err != nil {
			return -1, fmt.Errorf("%s: %v", op, err)
		}
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	if n == 0 {
		return n, fmt.Errorf("%s: no task with id: %v", op, id)
	}

	return n, nil
}

func (s *Storage) Delete(id int) (int64, error) {
	const op = "database.postgres.DeleteTodo"

	stmt, err := s.db.Prepare(`
	DELETE FROM public.todos 
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
		return n, fmt.Errorf("%s: no task with id: %v", op, id)
	}

	return n, nil
}

func (s *Storage) GetTodo(id int) (t.Todo, error) {
	const op = "database.postgres.GetTodo"

	rows, err := s.db.Query(`SELECT * FROM public.todos WHERE id = $1`, id)
	if err != nil {
		return t.Todo{}, fmt.Errorf("%s: %v", op, err)
	}

	var todo t.Todo

	if rows.Next() {
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Created, &todo.IsDone); err != nil {
			return t.Todo{}, fmt.Errorf("%s: %v", op, err)
		}
	} else {
		return t.Todo{}, fmt.Errorf("%s: no such task", op)
	}

	return todo, nil
}

func (s *Storage) OutputAll(filter string) ([]t.Todo, t.TodoInfo, int, error) {
	const op = "database.postgres.OutputAllTodos"

	query := ``
	switch filter {
	case "all":
		query = `SELECT * FROM public.todos`
	case "completed":
		query = `SELECT * FROM public.todos WHERE is_done = true`
	case "inWork":
		query = `SELECT * FROM public.todos WHERE is_done = false`
	default:
		query = `SELECT * FROM public.todos`
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, t.TodoInfo{}, 0, fmt.Errorf("%s: %v", op, err)
	}

	var result []t.Todo
	var todo t.Todo

	for rows.Next() {
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Created, &todo.IsDone); err != nil {
			return nil, t.TodoInfo{}, 0, fmt.Errorf("%s: %v", op, err)
		}

		result = append(result, todo)
	}

	var info t.TodoInfo

	query = `SELECT is_done FROM public.todos`

	rows, err = s.db.Query(query)
	if err != nil {
		return nil, t.TodoInfo{}, 0, fmt.Errorf("%s: %v", op, err)
	}

	var done bool
	for rows.Next() {
		if err := rows.Scan(&done); err != nil {
			return nil, t.TodoInfo{}, 0, fmt.Errorf("%s: %v", op, err)
		}

		if done {
			info.Completed++
		} else {
			info.InWork++
		}
		info.All++
	}

	return result, info, info.All, nil
}
