// This package provides CRUD operations with User
package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"

	util "github.com/sabbatD/srest-api/internal/http-server/handleUtil"
	"github.com/sabbatD/srest-api/internal/lib/api/access"
	resp "github.com/sabbatD/srest-api/internal/lib/api/response"
	"github.com/sabbatD/srest-api/internal/lib/logger/sl"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
)

type AuthResponse struct {
	resp.Response
	Token string `json:"token,omitempty"`
}

type GetResponse struct {
	resp.Response
	User u.TableUser `json:"user,omitempty"`
}

type UserHandler interface {
	Add(u u.User) (int, error)
	Auth(u u.AuthData) (succsess, isAdmin bool, id int, err error)
	Get(id int) (u.TableUser, error)
	UpdateUser(u u.User, id int) (int64, error)
}

// Register godoc
// @Summary Register a new user
// @Description Handles the registration of a new user by accepting a JSON payload containing user data.
// This endpoint will create a new user if the username doesn't already exist in the system.
// @Tags user
// @Accept json
// @Produce json
// @Param UserData body u.User true "Complete user data for registration"
// @Success 200 {object} GetResponse "Registration successful. Returns user data."
// @Failure 400 {object} resp.Response "Invalid input. Returns error message for improper data structure."
// @Router /signup [post]
func Register(log *slog.Logger, user UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.Register"

		log.With(util.SlogWith(op, r)...)

		var req u.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded")
		log.Debug("req: ", slog.Any("request", req))

		id, err := user.Add(req)
		if err != nil {
			if err.Error() == "database.postgres.Add: user already exists" {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error("user already exists"))

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		user, err := user.Get(id)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("user successfully created")
		log.Debug(fmt.Sprintf("user: %v", user))
		render.JSON(w, r, GetResponse{resp.OK(), user})
	}
}

// Auth godoc
// @Summary Authenticate user
// @Description Authenticates a user by accepting their login credentials (login and password) in JSON format.
// Upon successful authentication, a JWT token will be generated and returned for subsequent API calls.
// @Tags user
// @Accept json
// @Produce json
// @Param AuthData body u.AuthData true "User login credentials"
// @Success 200 {object} AuthResponse "Authentication successful. Returns a JWT token."
// @Failure 400 {object} resp.Response "Invalid input. Returns error message for improper data structure."
// @Router /signin [post]
func Auth(log *slog.Logger, user UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.Auth"

		log.With(util.SlogWith(op, r)...)

		var req u.AuthData
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded")
		log.Debug("req: ", slog.Any("request", req))

		ok, isAdmin, id, err := user.Auth(req)
		if err != nil {
			if !ok {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error("wrong login or password"))

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		token, err := access.GenerateJWT(id, req.Login, isAdmin)
		if err != nil {
			util.InternalError(w, r, log, fmt.Errorf("could not generate JWT Token"))
			return
		}

		log.Info("successfully logged in")
		log.Debug(fmt.Sprintf("user: %v", req))
		render.JSON(w, r, AuthResponse{resp.OK(), token})
	}
}

// Profile godoc
// @Summary Get user profile
// @Description Retrieves the full profile of the currently authenticated user.
// The user must be logged in and provide a valid JWT token for authentication.
// @Tags user
// @Produce json
// @Success 200 {object} GetResponse "Returns the user profile data."
// @Failure 400 {object} resp.Response "Invalid input. Returns error message for improper data structure."
// @Router /user/profile [get]
func Profile(log *slog.Logger, user UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.Profile"

		log.With(util.SlogWith(op, r)...)

		userContext, ok := r.Context().Value(access.CxtKey("userContext")).(access.UserContext)
		if !ok {
			http.Error(w, "User context not found", http.StatusUnauthorized)
			return
		}

		user, err := user.Get(userContext.Id)
		if err != nil {
			if err.Error() == "database.postgres.Get: no such user" {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error("No such user"))

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("User successfully retrieved")
		log.Debug(fmt.Sprintf("user: %v", user))
		render.JSON(w, r, GetResponse{resp.OK(), user})
	}
}

// UpdateUser godoc
// @Summary Update user profile
// @Description Updates the user profile with new data provided in the JSON payload.
// The user must be authenticated and provide a valid JWT token.
// @Tags user
// @Accept json
// @Produce json
// @Param Userdata body u.User true "Updated user data"
// @Success 200 {object} GetResponse "Profile successfully updated."
// @Failure 400 {object} resp.Response "Invalid input. Returns error message for improper data structure."
// @Router /user/profile [put]
func UpdateUser(log *slog.Logger, user UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.UpdateUser"

		log.With(util.SlogWith(op, r)...)

		var req u.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded")
		log.Debug("req: ", slog.Any("request", req))

		userContext, ok := r.Context().Value(access.CxtKey("userContext")).(access.UserContext)
		if !ok {
			http.Error(w, "User context not found", http.StatusUnauthorized)
			return
		}

		n, err := user.UpdateUser(req, userContext.Id)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))
			}
			util.InternalError(w, r, log, err)
			return
		}

		user, err := user.Get(userContext.Id)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("Successfully updated user")
		log.Debug(fmt.Sprintf("user: %v to %v with password %v and email %v", userContext.Id, req.Username, req.Password, req.Email))
		render.JSON(w, r, GetResponse{resp.OK(), user})
	}
}
