package user

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/sabbatD/srest-api/internal/lib/api/access"
	resp "github.com/sabbatD/srest-api/internal/lib/api/response"
	"github.com/sabbatD/srest-api/internal/lib/logger/sl"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
)

type RegisterResponse struct {
	resp.Response
	Authdata u.AuthData `json:"authdata,omitempty"`
}

type UserHandler interface {
	Add(u u.User) (int64, error)
	Auth(u u.AuthData) (succsess, isAdmin bool, err error)
	Get(username string) (u.TableUser, error)
	UpdateUser(u u.User, username string) (int64, error)
}

type GetResponse struct {
	resp.Response
	User u.TableUser `json:"user,omitempty"`
}

func Register(log *slog.Logger, user UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.user.Register"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req u.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		_, err := user.Add(req)
		if err != nil {
			if err.Error() == "database.postgres.Add: user already exists" {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error("user already exists"))

				return
			}
			log.Debug(err.Error())

			render.JSON(w, r, resp.Error("Internal Server Error"))

			return
		}

		log.Info("user successfully created")
		render.JSON(w, r, RegisterResponse{resp.OK(), u.AuthData{Username: req.Username, Password: req.Password}})
	}
}

func Auth(log *slog.Logger, user UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.user.Auth"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req u.AuthData
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		ok, isAdmin, err := user.Auth(req)
		if err != nil {
			if !ok {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error("wrong login or password"))

				return
			}
			log.Debug(err.Error())

			render.JSON(w, r, resp.Error("Internal Server Error"))

			return
		}

		token, err := access.GenerateJWT(req.Username, isAdmin)
		if err != nil {
			log.Debug("could not generate JWT Token")

			render.JSON(w, r, resp.Error("Internal Server Error"))

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(12 * time.Hour),
			HttpOnly: true,
		})

		log.Info("successfully logged in")
		render.JSON(w, r, resp.OK())
	}
}

func UpdateUser(log *slog.Logger, user UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.user.UpdateUser"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req u.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		userContext, ok := r.Context().Value("userContext").(access.UserContext)
		if !ok {
			http.Error(w, "User context not found", http.StatusUnauthorized)
			return
		}

		n, err := user.UpdateUser(req, userContext.Username)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))
			}
			log.Debug(err.Error())

			render.JSON(w, r, resp.Error("Internal Server Error"))

			return
		}

		log.Info(fmt.Sprintf("Successfully updated user: %v to %v with password %v and email %v", userContext.Username, req.Username, req.Password, req.Email))

		render.JSON(w, r, resp.OK())
	}

}

func Profile(log *slog.Logger, user UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.user.Profile"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userContext, ok := r.Context().Value("userContext").(access.UserContext)
		if !ok {
			http.Error(w, "User context not found", http.StatusUnauthorized)
			return
		}

		user, err := user.Get(userContext.Username)
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
