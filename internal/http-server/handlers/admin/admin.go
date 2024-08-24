package admin

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	access "github.com/sabbatD/srest-api/internal/lib/api/jwt"
	resp "github.com/sabbatD/srest-api/internal/lib/api/response"
	"github.com/sabbatD/srest-api/internal/lib/logger/sl"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
)

type GetAllResponse struct {
	resp.Response
	Users []u.TableUser `json:"users,omitempty"`
}

type GetResponse struct {
	resp.Response
	User u.TableUser `json:"user,omitempty"`
}

type AdminHandler interface {
	UpdateField(field string, u u.Login, val any) (int64, error)
	GetAll() ([]u.TableUser, error)
	Remove(u u.Login) (int64, error)
	Get(username string) (u.TableUser, error)
	UpdateUser(u u.User, username string) (int64, error)
}

func Update(log *slog.Logger, user AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.admin.Update"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		AdmCheck(w, r, log)

		var req u.Login
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			return
		}

		field := chi.URLParam(r, "field")
		var isIt bool

		switch field {
		case "block":
			isIt = true
		case "unblock":
			isIt = false
		case "makeadmin":
			isIt = true
		case "makeuser":
			isIt = false
		default:
			isIt = false
		}

		switch field {
		case "block":
			field = "block"
		case "unblock":
			field = "block"
		case "makeadmin":
			field = "admin"
		case "makeuser":
			field = "admin"
		default:
			log.Info("Wrong field in url: " + field)

			render.JSON(w, r, resp.Error("Wrong field in url"))

			return
		}

		if n, err := user.UpdateField(field, req, isIt); err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))
			}
			log.Debug(err.Error())

			render.JSON(w, r, resp.Error("Internal Server Error"))

			return
		}

		log.Info(fmt.Sprintf("Successfully updated field: %v to %v", field, isIt))

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

		AdmCheck(w, r, log)

		var req u.Login
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		n, err := user.Remove(req)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))

				return
			}
			log.Debug(err.Error())

			render.JSON(w, r, resp.Error("Internal Server Error"))

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

		AdmCheck(w, r, log)

		users, err := user.GetAll()
		if err != nil {
			if users == nil {
				log.Info("no users found")

				render.JSON(w, r, resp.Error("no users found"))

				return
			}
			log.Debug(err.Error())

			render.JSON(w, r, resp.Error("Internal Server Error"))

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

		AdmCheck(w, r, log)

		username := chi.URLParam(r, "username")

		user, err := user.Get(username)
		if err != nil {
			if err.Error() == "database.postgres.Get: no such user" {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error("No such user"))

				return
			}
			log.Debug(err.Error())

			render.JSON(w, r, resp.Error("Internal Server Error"))

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

		AdmCheck(w, r, log)

		var req u.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		username := chi.URLParam(r, "username")

		n, err := user.UpdateUser(req, username)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))
			}
			log.Debug(err.Error())

			render.JSON(w, r, resp.Error("Internal Server Error"))

			return
		}

		log.Info(fmt.Sprintf("Successfully updated user: %v to %v with password %v and email %v", username, req.Username, req.Password, req.Email))

		render.JSON(w, r, resp.OK())
	}
}

func contextAdmin(w http.ResponseWriter, r *http.Request) (bool, error) {
	userContext, ok := r.Context().Value("userContext").(access.UserContext)
	if !ok {
		http.Error(w, "User context not found", http.StatusUnauthorized)
		return false, fmt.Errorf("Unauthorized")
	}
	return userContext.IsAdmin, nil
}

func AdmCheck(w http.ResponseWriter, r *http.Request, log *slog.Logger) {
	ok, err := contextAdmin(w, r)
	if !ok {
		if err != nil {
			log.Error(err.Error())

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		log.Info("Not enough rights")

		render.JSON(w, r, resp.Error("Not enough rights"))

		return
	}
}
