package todo

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	util "github.com/sabbatD/srest-api/internal/http-server/handleUtil"
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
// @Description Handles the creation of a new task by accepting a JSON payload containing task data.
// @Tags todo
// @Accept json
// @Produce json
// @Param UserData body t.TodoRequest true "Complete task data for creation"
// @Success 200 {object}  t.Todo "Creation successful. Returns task with status code OK."
// @Failure 400 {object} string "failed to deserialize json request."
// @Failure 500 {object} string "Internal error."
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

// Get All godoc
// @Summary Get all tasks
// @Description Gets all tasks and returns a JSON containing task data.
// @Tags todo
// @Produce json
// @Param filter query string false "all, completed, or inWork"
// @Success 200 {object} t.MetaResponse "Retrieved successfully. Returns status code OK."
// @Failure 500 {object} string "Internal error."
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
// @Summary Get task
// @Description Gets a task by ID in the URL and returns a JSON containing task data.
// @Tags todo
// @Produce json
// @Success 200 {object}  t.Todo "Retrieved successfully. Returns task and status code OK."
// @Failure 400 {object} string "Missing or wrong id."
// @Failure 404 {object} string "No such task."
// @Failure 500 {object} string "Internal error."
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
// @Summary Update task
// @Description Handles the update of a task by accepting a JSON payload containing task data.
// @Tags todo
// @Accept json
// @Produce json
// @Param UserData body t.TodoRequest true "Complete task data for update"
// @Success 200 {object}  t.Todo "Update successful. Returns task with status code OK."
// @Failure 400 {object} string "failed to deserialize json request."
// @Failure 400 {object} string "Missing or wrong id."
// @Failure 404 {object} string "No such task."
// @Failure 500 {object} string "Internal error."
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
// @Summary Delete task
// @Description Deletes a task by ID in the URL.
// @Tags todo
// @Produce json
// @Success 200 {object} string "Deletion successful. Returns status code OK."
// @Failure 400 {object} string "Missing or wrong id."
// @Failure 404 {object} string "No such task."
// @Failure 500 {object} string "Internal error."
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
			}
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("successfully retrieved task")
	}
}
