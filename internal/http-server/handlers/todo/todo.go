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
	Todos t.Todos
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
	OutputAll(isDone bool) ([]t.Todo, error)
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
				render.JSON(w, r, resp.Error(err.Error()))
				return
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info("successfully retrieved task")
		render.JSON(w, r, GetResponse{resp.OK(), task})
	}
}
func GetAll(log *slog.Logger, todo TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.todo.GetAll"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		filterstr := r.URL.Query().Get("filter")
		filter, _ := strconv.ParseBool(filterstr)

		todos, err := todo.OutputAll(filter)
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		log.Info("successfully retrieved task")
		render.JSON(w, r, GetAllResponse{resp.OK(), todos})
	}
}

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
