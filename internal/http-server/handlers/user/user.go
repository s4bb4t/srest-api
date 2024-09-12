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

type AccessToken struct {
	Token string `json:"access"`
}
type RefreshToken struct {
	Token string `json:"refresh"`
}

type Tokens struct {
	AccessToken
	RefreshToken
}

type AuthResponse struct {
	resp.Response
	Tokens `json:"tokens,omitempty"`
}

type GetResponse struct {
	resp.Response
	User u.TableUser `json:"user,omitempty"`
}

type UserHandler interface {
	Add(u u.User) (int, error)
	Auth(u u.AuthData) (user u.TableUser, err error)
	Get(id int) (u.TableUser, error)
	UpdateUser(u u.User, id int) (int64, error)
	RefreshToken(token string) (string, int, error)
	SaveRefreshToken(token string, id int) error
}

// Register godoc
// @Summary Register a new user
// @Description Handles the registration of a new user by accepting a JSON payload containing user data.
// This endpoint will create a new user if the username doesn't already exist in the system.
// @Tags user
// @Accept json
// @Produce json
// @Param UserData body u.User true "Complete user data for registration"
// @Success 201 {object} GetResponse "Registration successful. Returns user data."
// @Failure 400 {object}  "Invalid input."
// @Failure 409 {object}  "User already exists."
// @Failure 500 {object}  "Internal error."
// @Router /signup [post]
func Register(log *slog.Logger, User UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.Register"

		log.With(util.SlogWith(op, r)...)

		var req u.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			http.Error(w, "failed to deserialize json request", http.StatusBadRequest)

			return
		}

		log.Info("request body decoded")
		log.Debug("req: ", slog.Any("request", req))

		id, err := User.Add(req)
		if err != nil {
			if err.Error() == "database.postgres.Add: user already exists" {
				log.Info(err.Error())

				http.Error(w, "user already exists", http.StatusConflict)

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		user, err := User.Get(id)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("user successfully created")
		log.Debug(fmt.Sprintf("user: %v", user))

		w.WriteHeader(http.StatusCreated)
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
// @Failure 400 {object}  "Invalid input."
// @Failure 401 {object}  "Invalid credentials."
// @Failure 500 {object}  "Internal error."
// @Router /signin [post]
func Auth(log *slog.Logger, User UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.Auth"

		log.With(util.SlogWith(op, r)...)

		var req u.AuthData
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			http.Error(w, "failed to deserialize json request", http.StatusBadRequest)

			return
		}

		log.Info("request body decoded")
		log.Debug("req: ", slog.Any("request", req))

		user, err := User.Auth(req)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}
		if user.ID == 0 {
			log.Info("wrong login or password")

			http.Error(w, "Invalid credentials", http.StatusUnauthorized)

			return
		}

		accessToken, err := access.NewAccessToken(user.ID, user.Login, user.Admin)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		refreshToken := access.NewRefreshToken()
		if refreshToken == "" {
			util.InternalError(w, r, log, fmt.Errorf("could not generate JWT refreshToken"))
			return
		}

		if err := User.SaveRefreshToken(refreshToken, user.ID); err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("successfully logged in")
		log.Debug(fmt.Sprintf("user: %v", req))

		render.JSON(w, r, AuthResponse{resp.OK(), Tokens{AccessToken{accessToken}, RefreshToken{refreshToken}}})
	}
}

// Refresh godoc
// @Summary Refresh user's access token
// @Description Recieve a user's refresh token in JSON format.
// Upon successful refresh token compare, an access JWT token will be generated and returned for subsequent API calls.
// @Tags user
// @Accept json
// @Produce json
// @Param AuthData body RefreshToken true "User's refresh token"
// @Success 200 {object} AuthResponse "Authentication successful. Returns a JWT token."
// @Failure 400 {object}  "Invalid input."
// @Failure 401 {object}  "Invalid credentials: token is expired - must auth again."
// @Failure 500 {object}  "Internal error."
// @Router /auth/refresh [post]
func Refresh(log *slog.Logger, User UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.Refresh"

		log.With(util.SlogWith(op, r)...)

		var req RefreshToken
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			http.Error(w, "failed to deserialize json request", http.StatusBadRequest)

			return
		}

		log.Info("request body decoded")
		log.Debug("req: ", slog.Any("request", req))

		token, id, err := User.RefreshToken(req.Token)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}
		if token == "expired" {
			log.Info("token is expired")

			http.Error(w, "Invalid credentials: token is expired - must auth again", http.StatusUnauthorized)

			return
		}

		user, err := User.Get(id)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		accessToken, err := access.NewAccessToken(user.ID, user.Login, user.Admin)
		if err != nil {
			util.InternalError(w, r, log, fmt.Errorf("could not generate JWT accessToken"))
			return
		}

		log.Info("successfully refreshed access token")
		log.Debug(fmt.Sprintf("user: %v", req))

		render.JSON(w, r, AuthResponse{resp.OK(), Tokens{AccessToken{accessToken}, RefreshToken{req.Token}}})
	}
}

// Profile godoc
// @Summary Get user profile
// @Description Retrieves the full profile of the currently authenticated user.
// The user must be logged in and provide a valid JWT token for authentication.
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} GetResponse "Returns the user profile data."
// @Failure 400 {object}  "No such user."
// @Failure 500 {object}  "Internal error."
// @Router /user/profile [get]
func Profile(log *slog.Logger, User UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.Profile"

		log.With(util.SlogWith(op, r)...)

		userContext, ok := r.Context().Value(access.CxtKey("userContext")).(access.UserContext)
		if !ok {
			http.Error(w, "User context not found", http.StatusUnauthorized)
			return
		}

		user, err := User.Get(userContext.Id)
		if err != nil {
			if err.Error() == "database.postgres.Get: no such user" {
				log.Info(err.Error())

				http.Error(w, "No such user", http.StatusBadRequest)

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
// @Security BearerAuth
// @Success 200 {object} GetResponse "Profile successfully updated."
// @Failure 400 {object}  "Invalid input."
// @Failure 400 {object}  "Login or email already used."
// @Failure 404 {object}  "No such user."
// @Failure 500 {object}  "Internal error."
// @Router /user/profile [put]
func UpdateUser(log *slog.Logger, User UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.UpdateUser"

		log.With(util.SlogWith(op, r)...)

		var req u.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			http.Error(w, "failed to deserialize json request", http.StatusBadRequest)

			return
		}

		log.Info("request body decoded")
		log.Debug("req: ", slog.Any("request", req))

		userContext, ok := r.Context().Value(access.CxtKey("userContext")).(access.UserContext)
		if !ok {
			http.Error(w, "User context not found", http.StatusUnauthorized)
			return
		}

		n, err := User.UpdateUser(req, userContext.Id)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				http.Error(w, "No such user", http.StatusNotFound)
			} else if n == -2 {

				log.Info(err.Error())

				http.Error(w, "Login or email already used", http.StatusBadRequest)
			}
			util.InternalError(w, r, log, err)
			return
		}

		user, err := User.Get(userContext.Id)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("Successfully updated user")
		log.Debug(fmt.Sprintf("user: %v to %v with password %v and email %v", userContext.Id, req.Username, req.Password, req.Email))

		render.JSON(w, r, GetResponse{resp.OK(), user})
	}
}
