package database

import (
	"database/sql"
	"fmt"
	"time"

	t "github.com/sabbatD/srest-api/internal/lib/todoConfig"
)

func (s *Storage) Create(t t.TodoRequest) error {
	const op = "database.postgres.CreateTodo"

	stmt, err := s.db.Prepare(`
		INSERT INTO public.todos (
			id, title, created
		) VALUES ($1, $2, $3)
	`)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	var maxID sql.NullInt64

	err = s.db.QueryRow(`SELECT MAX(id) FROM public.todos;`).Scan(&maxID)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if !maxID.Valid {
		maxID.Int64 = 1
	} else {
		maxID.Int64 += 1
	}

	_, err = stmt.Exec(maxID.Int64, t.Title, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (s *Storage) Update(id int, t t.TodoRequest) (int64, error) {
	const op = "database.postgres.UpdateTodo"

	stmt, err := s.db.Prepare(`
		UPDATE public.todos 
			SET title = $1
			WHERE id = $2
	`)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.Exec(t.Title, id)
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

func (s *Storage) OutputAll(isDone bool) ([]t.Todo, error) {
	const op = "database.postgres.OutputAllTodos"

	query := `SELECT * FROM public.todos WHERE isDone = $1`

	rows, err := s.db.Query(query, isDone)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	var result []t.Todo
	var todo t.Todo

	for rows.Next() {
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Created, &todo.IsDone); err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}

		result = append(result, todo)
	}

	return result, nil
}
