// Package todo provides handlers for managing tasks in a TODO application.
// It supports operations such as creating, updating, retrieving, and deleting tasks.
// The handlers accept and return JSON data, and include support for filtering tasks based on their status.
package todo

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	util "github.com/sabbatD/srest-api/internal/http-server/handleUtil"
	"github.com/sabbatD/srest-api/internal/lib/api/validation"
	"github.com/sabbatD/srest-api/internal/lib/logger/sl"
	t "github.com/sabbatD/srest-api/internal/lib/todoConfig"
)

type TodoHandler interface {
	Create(t t.TodoRequest) (int64, error)
	Update(id int, t t.TodoRequest) (int64, error)
	Delete(id int) (int64, error)
	GetTodo(id int) (t.Todo, error)
	OutputAll(filter string) ([]t.Todo, t.TodoInfo, int, error)
}

// Create godoc
// @Summary Create a new task
// @Description Creates a new task by accepting a JSON payload with the task's details.
// @Tags todo
// @Accept json
// @Produce json
// @Param UserData body t.TodoRequest true "Task data for creating a new task"
// @Success 200 {object}  t.Todo "Task successfully created, returns the created task."
// @Failure 400 {object} string "Invalid request body or missing/incorrect fields."
// @Failure 500 {object} string "Internal server error."
// @Router /todos [post]
func Create(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.Create"

		log.With(util.SlogWith(op, r)...)

		var req t.TodoRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			http.Error(w, "failed to deserialize json request", http.StatusBadRequest)

			return
		}

		log.Info("request body decoded")
		log.Debug("req: ", slog.Any("request", req))

		id, err := todo.Create(req)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		task, err := todo.GetTodo(int(id))
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("successfully created task")

		render.JSON(w, r, task)
	}
}

// GetAll godoc
// @Summary Retrieve all tasks
// @Description Retrieves all tasks with optional filtering by status (e.g., completed or in-progress).
// @Tags todo
// @Produce json
// @Param filter query string false "Filter tasks by status: all, completed, or inWork"
// @Success 200 {object} t.MetaResponse "Tasks retrieved successfully."
// @Failure 500 {object} string "Internal server error."
// @Router /todos [get]
func GetAll(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.GetAll"

		log.With(util.SlogWith(op, r)...)

		filter := r.URL.Query().Get("filter")

		todos, info, n, err := todo.OutputAll(filter)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		if todos == nil {
			todos = []t.Todo{}
		}

		response := t.MetaResponse{
			Data: todos,
			Info: info,
			Meta: t.Meta{
				TotalAmount: n,
			},
		}

		log.Info("successfully retrieved tasks")

		render.JSON(w, r, response)
	}
}

// Get godoc
// @Summary Retrieve a task by ID
// @Description Retrieves a specific task by its ID from the URL.
// @Tags todo
// @Produce json
// @Param id path int true "ID of the task to retrieve"
// @Success 200 {object}  t.Todo "Task retrieved successfully."
// @Failure 400 {object} string "Invalid or missing task ID."
// @Failure 404 {object} string "Task not found."
// @Failure 500 {object} string "Internal server error."
// @Router /todos/{id} [get]
func Get(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.Get"

		log.With(util.SlogWith(op, r)...)

		id := util.GetUrlParam(w, r, log)
		if id == 0 {
			log.Info("missing or wrong id")
			http.Error(w, "Missing or wrong id", http.StatusBadRequest)
			return
		}

		task, err := todo.GetTodo(id)
		if err != nil {
			if err.Error() == "database.postgres.GetTodo: no such task" {
				log.Info(err.Error())

				http.Error(w, "No such task", http.StatusNotFound)

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("successfully retrieved task")

		render.JSON(w, r, task)
	}
}

// Update godoc
// @Summary Update an existing task
// @Description Updates an existing task by accepting a JSON payload with the updated task details.
// @Tags todo
// @Accept json
// @Produce json
// @Param id path int true "ID of the task to update"
// @Param UserData body t.TodoRequest true "Updated task data"
// @Success 200 {object}  t.Todo "Task updated successfully, returns the updated task."
// @Failure 400 {object} string "Invalid request body, missing/incorrect fields, or invalid ID."
// @Failure 404 {object} string "Task not found."
// @Failure 500 {object} string "Internal server error."
// @Router /todos/{id} [put]
func Update(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.Update"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req t.TodoRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			http.Error(w, "failed to deserialize json request", http.StatusBadRequest)

			return
		}

		log.Info("request body decoded")
		log.Debug("req: ", slog.Any("request", req))

		validation.InitValidator()
		if err := validation.ValidateStruct(req); err != nil {
			log.Debug(fmt.Sprintf("validation failed: %v", err.Error()))

			http.Error(w, fmt.Sprintf("Invalid input: %v", err.Error()), http.StatusBadRequest)

			return
		}

		log.Info("input validated")

		id := util.GetUrlParam(w, r, log)
		if id == 0 {
			log.Info("missing or wrong id")
			http.Error(w, "Missing or wrong id", http.StatusBadRequest)
			return
		}

		n, err := todo.Update(id, req)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				http.Error(w, "No such task", http.StatusNotFound)

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		task, err := todo.GetTodo(id)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("successfully updated task")

		render.JSON(w, r, task)
	}
}

// Delete godoc
// @Summary Delete a task by ID
// @Description Deletes a task by its ID from the URL.
// @Tags todo
// @Produce json
// @Param id path int true "ID of the task to delete"
// @Success 200 {object} string "Task deleted successfully."
// @Failure 400 {object} string "Invalid or missing task ID."
// @Failure 404 {object} string "Task not found."
// @Failure 500 {object} string "Internal server error."
// @Router /todos/{id} [delete]
func Delete(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.Delete"

		log.With(util.SlogWith(op, r)...)

		id := util.GetUrlParam(w, r, log)
		if id == 0 {
			log.Info("missing or wrong id")

			http.Error(w, "Missing or wrong id", http.StatusBadRequest)

			return
		}

		n, err := todo.Delete(id)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				http.Error(w, "No such task", http.StatusNotFound)

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("successfully retrieved task")
	}
}
