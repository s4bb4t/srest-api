package todo

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/sabbatD/srest-api/internal/lib/api/response"
	"github.com/sabbatD/srest-api/internal/lib/logger/sl"
	t "github.com/sabbatD/srest-api/internal/lib/todoConfig"
)

type GetAllResponse struct {
	resp.Response
	t.MetaResponse `json:"metaresponse,omitempty"`
}

type GetResponse struct {
	resp.Response
	Todo t.Todo `json:"todo,omitempty"`
}

type TodoHandler interface {
	Create(t t.TodoRequest) error
	Update(id int, t t.TodoRequest) (int64, error)
	Delete(id int) (int64, error)
	GetTodo(id int) (t.Todo, error)
	OutputAll(filter string) ([]t.Todo, t.TodoInfo, int, error)
}

// InternalError is a shortcut for internal error handling
// InternalError always returns an error to debug logs and JSON ERROR
func InternalError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.Debug(err.Error())

	render.JSON(w, r, resp.Error("Internal Server Error"))
}

// JsonDecodeError is a shortcut for json.UnmarshalError
// InternalError always returns an error and JSON ERROR
func JsonDecodeError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.Error("failed to decode request body", sl.Err(err))

	render.JSON(w, r, resp.Error("failed to decode request"))
}

// Register godoc
// @Summary Create a new task
// @Description Handles the creation of a new task by accepting a JSON payload containing task data.
// @Tags user
// @Accept json
// @Produce json
// @Param UserData body t.TodoRequest true "Complete task data for creation"
// @Success 200 {object} resp.Response "Creation successful. Returns status code OK."
// @Failure 401 {object} resp.Response "Creation failed. Returns error message."
// @Router /todos [post]
func Create(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.Create"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req t.TodoRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			JsonDecodeError(w, r, log, err)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err := todo.Create(req)
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		log.Info("successfully created task")
		render.JSON(w, r, resp.OK())
	}
}

// Register godoc
// @Summary Update task
// @Description Handles the upd of a task by accepting a JSON payload containing task data.
// @Tags user
// @Accept json
// @Produce json
// @Param UserData body t.TodoRequest true "Complete task data for creation"
// @Success 200 {object} resp.Response "Creation successful. Returns status code OK."
// @Failure 401 {object} resp.Response "Creation failed. Returns error message."
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
			JsonDecodeError(w, r, log, err)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		n, err := todo.Update(id, req)
		if err != nil {
			if n == 0 {
				render.JSON(w, r, resp.Error(err.Error()))
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info("successfully updated task")
		render.JSON(w, r, resp.OK())
	}
}

// Register godoc
// @Summary Get task
// @Description Gets a task by id in url and return a JSON containing task data.
// @Tags user
// @Produce json
// @Success 200 {object} GetResponse "Retrieved successfully. Returns status code OK."
// @Failure 401 {object} resp.Response "Retrieving failed. Returns error message."
// @Router /todos/{id} [get]
func Get(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.Get"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		task, err := todo.GetTodo(id)
		if err != nil {
			if err.Error() == "database.postgres.GetTodo: no such task" {
				render.JSON(w, r, resp.Error("no such task"))
				return
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info("successfully retrieved task")
		render.JSON(w, r, GetResponse{resp.OK(), task})
	}
}

// Register godoc
// @Summary Get all tasks
// @Description Gets all tasks and return a JSON containing tasks data.
// @Tags user
// @Produce json
// @Param filter query string true "all, completed or inwork"
// @Success 200 {object} GetAllResponse "Retrieved successfully. Returns status code OK."
// @Failure 401 {object} resp.Response "Retrieving failed. Returns error message."
// @Router /todos [get]
func GetAll(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.GetAll"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		filter := r.URL.Query().Get("filter")

		todos, info, n, err := todo.OutputAll(filter)
		if err != nil {
			if err.Error() == "database.postgres.OutputAllTodos: unknown filter" {
				render.JSON(w, r, resp.Error("unknown filter"))
			}
			InternalError(w, r, log, err)
			return
		}

		response := GetAllResponse{
			resp.OK(),
			t.MetaResponse{
				Data: todos,
				Info: info,
				Meta: t.Meta{
					TotalAmount: n,
				},
			},
		}

		log.Info("successfully retrieved task")
		render.JSON(w, r, response)
	}
}

// Register godoc
// @Summary Delete task
// @Description Delete task by id in url.
// @Tags user
// @Produce json
// @Success 200 {object} resp.Response "Creation successful. Returns status code OK."
// @Failure 401 {object} resp.Response "Creation failed. Returns error message."
// @Router /todos/{id} [delete]
func Delete(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.Delete"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		n, err := todo.Delete(id)
		if err != nil {
			if n == 0 {
				render.JSON(w, r, resp.Error(err.Error()))
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info("successfully retrieved task")
		render.JSON(w, r, resp.OK())
	}
}
