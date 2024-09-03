package admin

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sabbatD/srest-api/internal/lib/api/access"
	resp "github.com/sabbatD/srest-api/internal/lib/api/response"
	"github.com/sabbatD/srest-api/internal/lib/logger/sl"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
)

type UpdateRequest struct {
	Field string
	Value any
}

type GetAllResponse struct {
	resp.Response
	Users []u.TableUser `json:"users,omitempty"`
}

type GetResponse struct {
	resp.Response
	User u.TableUser `json:"user,omitempty"`
}

type AdminHandler interface {
	UpdateField(field string, id int, val any) (int64, error)
	GetAll(search, order string, blocked bool, limit, offset int) ([]u.TableUser, error)
	Remove(id int) (int64, error)
	Get(id int) (u.TableUser, error)
	UpdateUser(u u.User, id int) (int64, error)
}

func Update(log *slog.Logger, user AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.admin.Update"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if !AdmCheck(w, r, log) {
			return
		}

		var req UpdateRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			JsonDecodeError(w, r, log, err)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		if n, err := user.UpdateField(req.Field, id, req.Value); err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info(fmt.Sprintf("Successfully updated field: %v to %v", req.Field, req.Value))

		render.JSON(w, r, resp.OK())
	}
}

func Block(log *slog.Logger, user AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.admin.Block"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if !AdmCheck(w, r, log) {
			return
		}

		var req UpdateRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			JsonDecodeError(w, r, log, err)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		if n, err := user.UpdateField("block", id, true); err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info(fmt.Sprintf("Successfully updated field: %v to %v", req.Field, req.Value))

		render.JSON(w, r, resp.OK())
	}
}

func Unblock(log *slog.Logger, user AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.admin.Block"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if !AdmCheck(w, r, log) {
			return
		}

		var req UpdateRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			JsonDecodeError(w, r, log, err)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		if n, err := user.UpdateField("block", id, false); err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info(fmt.Sprintf("Successfully updated field: %v to %v", req.Field, req.Value))

		render.JSON(w, r, resp.OK())
	}
}

func Remove(log *slog.Logger, user AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.admin.Remove"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if !AdmCheck(w, r, log) {
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		n, err := user.Remove(id)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))

				return
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info("user successfully removed")

		render.JSON(w, r, resp.OK())
	}
}

func GetAll(log *slog.Logger, user AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.admin.GetAll"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if !AdmCheck(w, r, log) {
			return
		}

		search := r.URL.Query().Get("search")
		order := r.URL.Query().Get("order")
		blockedStr := r.URL.Query().Get("blocked")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		if search == "" || order == "" || blockedStr == "" || limitStr == "" || offsetStr == "" {
			log.Info("one or more parameters are empty")

			render.JSON(w, r, resp.Error("one or more parameters are empty"))

			return
		}

		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)
		blocked, _ := strconv.ParseBool(offsetStr)

		users, err := user.GetAll(search, order, blocked, limit, offset)
		if err != nil {
			if users == nil {
				log.Info("no users found")

				render.JSON(w, r, resp.Error("no users found"))

				return
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info("users successfully retrieved")

		render.JSON(w, r, GetAllResponse{resp.OK(), users})
	}
}

func Profile(log *slog.Logger, user AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.user.Profile"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if !AdmCheck(w, r, log) {
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		user, err := user.Get(id)
		if err != nil {
			if err.Error() == "database.postgres.Get: no such user" {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error("No such user"))

				return
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info("user successfully retrieved")
		render.JSON(w, r, GetResponse{resp.OK(), user})
	}
}

func UpdateUser(log *slog.Logger, user AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.admin.UpdateUser"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if !AdmCheck(w, r, log) {
			return
		}

		var req u.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			InternalError(w, r, log, err)
			return
		}

		n, err := user.UpdateUser(req, id)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))
			}
			InternalError(w, r, log, err)
			return
		}

		log.Info(fmt.Sprintf("Successfully updated user: %v to %v with password %v and email %v", id, req.Username, req.Password, req.Email))

		render.JSON(w, r, resp.OK())
	}
}

func contextAdmin(w http.ResponseWriter, r *http.Request) (bool, error) {
	userContext, ok := r.Context().Value(access.CxtKey("userContext")).(access.UserContext)
	if !ok {
		http.Error(w, "User context not found", http.StatusUnauthorized)
		return false, fmt.Errorf("Unauthorized")
	}
	return userContext.IsAdmin, nil
}

func AdmCheck(w http.ResponseWriter, r *http.Request, log *slog.Logger) bool {
	ok, err := contextAdmin(w, r)
	if !ok {
		if err != nil {
			log.Error(err.Error())

			render.JSON(w, r, resp.Error(err.Error()))

			return false
		}

		log.Info("Not enough rights")

		render.JSON(w, r, resp.Error("Not enough rights"))

		return false
	}
	return true
}

func InternalError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.Debug(err.Error())

	render.JSON(w, r, resp.Error("Internal Server Error"))
}

func JsonDecodeError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.Error("failed to decode request body", sl.Err(err))

	render.JSON(w, r, resp.Error("failed to decode request"))
}
