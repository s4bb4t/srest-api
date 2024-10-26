package admin

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	util "github.com/sabbatD/srest-api/internal/http-server/handleUtil"
	"github.com/sabbatD/srest-api/internal/lib/api/access"
	"github.com/sabbatD/srest-api/internal/lib/api/validation"
	"github.com/sabbatD/srest-api/internal/lib/logger/sl"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
)

type UpdateRequest struct {
	Field string
	Value any
}

type AdminHandler interface {
	UpdateField(field string, id int, val any) (int64, error)
	All(q u.GetAllQuery) (result u.MetaResponse, E error)
	Remove(id int) (int64, error)
	Get(id int) (u.TableUser, error)
	UpdateUser(u u.PutUser, id int) (int64, error)
}

// All godoc
// @Summary Get all users
// @Description Fetches a list of users based on optional query parameters such as filters and sorting.
// Requires Authorization header with Bearer token for authentication.
// @Tags admin
// @Produce json
// @Param search query string false "Filter users by username or email"
// @Param sortBy query string false "Sort by 'email', 'username', or 'id'. Default is 'id'."
// @Param sortOrder query string false "Sort order: 'asc', 'desc', or 'none'. Default is 'asc'."
// @Param isBlocked query bool false "Filter by block status (true/false)"
// @Param limit query int false "Limit the number of users returned (default is 20)"
// @Param offset query int false "Offset for pagination (default is 0)"
// @Security BearerAuth
// @Success 200 {object} u.MetaResponse "Successful retrieval of users."
// @Failure 401 {object} string "Unauthorized access. Bearer token missing or invalid."
// @Failure 403 {object} string "Insufficient permissions."
// @Failure 500 {object} string "Internal server error."
// @Router /admin/users [get]
func All(log *slog.Logger, Users AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.GetAll"

		log.With(util.SlogWith(op, r)...)

		if !AdmCheck(w, r, log) {
			return
		}

		var q u.GetAllQuery
		var E error

		q.SearchTerm = r.URL.Query().Get("search")

		q.SortBy = strings.ToLower(r.URL.Query().Get("sortBy"))
		switch q.SortBy {
		case "email", "username", "id":
		default:
			q.SortBy = "id"
		}

		q.SortOrder = strings.ToUpper(r.URL.Query().Get("sortOrder"))
		switch q.SortOrder {
		case "ASC", "DESC", "NONE":
		default:
			q.SortOrder = "ASC"
		}

		isBlockedStr := r.URL.Query().Get("isBlocked")
		q.IsBlocked, E = strconv.ParseBool(isBlockedStr)
		if E != nil {
			q.IsBlocked = false
		}

		limitStr := r.URL.Query().Get("limit")
		q.Limit, E = strconv.Atoi(limitStr)
		if E != nil || q.Limit < 20 {
			q.Limit = 20
		}

		fmt.Println("http", limitStr, q.Limit, E)

		offsetStr := r.URL.Query().Get("offset")
		q.Offset, E = strconv.Atoi(offsetStr)
		if E != nil || q.Offset < 0 {
			q.Offset = 0
		}

		metaResponse, err := Users.All(q)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("users successfully retrieved")
		log.Debug(fmt.Sprintf("query: %v", q))

		render.JSON(w, r, metaResponse)
	}
}

// Profile godoc
// @Summary Retrieve user's profile
// @Description Retrieves a user's profile by their ID.
// Requires Authorization header with Bearer token for authentication.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID of the user"
// @Success 200 {object} u.TableUser "Successful retrieval of user profile."
// @Failure 400 {object} string "Invalid or missing user ID."
// @Failure 401 {object} string "Unauthorized access. Bearer token missing or invalid."
// @Failure 403 {object} string "Insufficient permissions."
// @Failure 404 {object} string "User not found."
// @Failure 500 {object} string "Internal server error."
// @Router /admin/users/{id} [get]
func Profile(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.user.Profile"

		log.With(util.SlogWith(op, r)...)

		if !AdmCheck(w, r, log) {
			return
		}

		id := util.GetUrlParam(w, r, log)
		if id == 0 {
			log.Info("missing or wrong id")
			http.Error(w, "Missing or wrong id", http.StatusBadRequest)
			return
		}

		user, err := User.Get(id)
		if err != nil {
			if err.Error() == "database.postgres.Get: no such user" {
				log.Info(err.Error())

				http.Error(w, "No such user", http.StatusNotFound)

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("user successfully retrieved")
		log.Debug(fmt.Sprintf("user: %v", user))

		render.JSON(w, r, user)
	}
}

// UpdateUser godoc
// @Summary Update user's profile
// @Description Updates the details of a user by accepting a JSON payload.
// Requires Authorization header with Bearer token for authentication.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "ID of the user"
// @Param UserData body u.PutUser true "User data payload"
// @Security BearerAuth
// @Success 200 {object} u.TableUser "User profile updated successfully."
// @Failure 400 {object} string "Invalid request payload or ID."
// @Failure 400 {object} string "Duplicate login or email."
// @Failure 401 {object} string "Unauthorized access. Bearer token missing or invalid."
// @Failure 403 {object} string "Insufficient permissions."
// @Failure 404 {object} string "User not found."
// @Failure 500 {object} string "Internal server error."
// @Router /admin/users/{id} [put]
func UpdateUser(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.UpdateUser"

		log.With(util.SlogWith(op, r)...)

		if !AdmCheck(w, r, log) {
			return
		}

		var req u.PutUser
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

		n, err := User.UpdateUser(req, id)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				http.Error(w, "No such user", http.StatusNotFound)

				return
			} else if n == -2 {

				log.Info(err.Error())

				http.Error(w, "Login or email already used", http.StatusBadRequest)

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		user, err := User.Get(id)
		if err != nil {
			util.InternalError(w, r, log, err)
		}

		log.Info("Successfully updated user")
		log.Debug(fmt.Sprintf("user: %v", user))

		render.JSON(w, r, user)
	}
}

// Remove godoc
// @Summary Remove user
// @Description Deletes a user by their ID.
// Requires Authorization header with Bearer token for authentication.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID of the user"
// @Success 200 {object} string "User successfully removed."
// @Failure 400 {object} string "Invalid or missing user ID."
// @Failure 401 {object} string "Unauthorized access. Bearer token missing or invalid."
// @Failure 403 {object} string "Insufficient permissions."
// @Failure 404 {object} string "User not found."
// @Failure 500 {object} string "Internal server error."
// @Router /admin/users/{id} [delete]
func Remove(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.Remove"

		log.With(util.SlogWith(op, r)...)

		if !AdmCheck(w, r, log) {
			return
		}

		id := util.GetUrlParam(w, r, log)
		if id == 0 {
			log.Info("missing or wrong id")
			http.Error(w, "Missing or wrong id", http.StatusBadRequest)
			return
		}

		n, err := User.Remove(id)
		if err != nil {
			if n == 0 {
				log.Info(err.Error())

				http.Error(w, "No such user", http.StatusNotFound)

				return
			}
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("user successfully removed")
	}
}

// Block godoc
// @Summary Block user
// @Description Blocks a user by their ID, disabling their account.
// Requires Authorization header with Bearer token for authentication.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID of the user"
// @Success 200 {object} u.TableUser "User successfully blocked."
// @Failure 400 {object} string "Invalid or missing user ID."
// @Failure 404 {object} string "User not found."
// @Failure 500 {object} string "Internal server error."
// @Router /admin/users/{id}/block [post]
func Block(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.Block"

		changeField(w, r, log, User, op, "block", true)
	}
}

// Unblock godoc
// @Summary Unlock user
// @Description Unblocks a user by their ID, re-enabling their account.
// Requires Authorization header with Bearer token for authentication.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID of the user"
// @Success 200 {object} u.TableUser "User successfully unblocked."
// @Failure 400 {object} string "Invalid or missing user ID."
// @Failure 404 {object} string "User not found."
// @Failure 500 {object} string "Internal server error."
// @Router /admin/users/{id}/unlock [post]
func Unblock(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.Unblock"

		changeField(w, r, log, User, op, "block", false)
	}
}

// Update godoc
// @Summary Update user's rights
// @Description Updates specific fields related to user's rights by accepting a JSON payload.
// Requires Authorization header with Bearer token for authentication.
// @Tags admin
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID of the user"
// @Param UserData body UpdateRequest true "User data for updating rights"
// @Success 200 {object} u.TableUser "Rights successfully updated."
// @Failure 400 {object} string "Invalid request payload or missing ID."
// @Failure 400 {object} string "No such field."
// @Failure 404 {object} string "User not found."
// @Failure 500 {object} string "Internal server error."
// @Router /admin/users/{id}/rights [post]
func Update(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.Update"

		log.With(util.SlogWith(op, r)...)

		if !AdmCheck(w, r, log) {
			return
		}

		var req UpdateRequest
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

		if n, err := User.UpdateField(req.Field, id, req.Value); err != nil {
			if n == 0 {
				log.Info(err.Error())

				http.Error(w, "No such user", http.StatusNotFound)

				return
			} else if n == -2 {
				log.Info(err.Error())

				http.Error(w, "No such field", http.StatusBadRequest)

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

		log.Info("Successfully updated user's rights")
		log.Debug(fmt.Sprintf("user: %v", user))

		render.JSON(w, r, user)
	}
}

func contextAdmin(r *http.Request) (bool, error) {
	userContext, ok := r.Context().Value(access.CxtKey("userContext")).(access.UserContext)
	if !ok {
		return false, fmt.Errorf("unauthorized")
	}
	return userContext.IsAdmin, nil
}

func AdmCheck(w http.ResponseWriter, r *http.Request, log *slog.Logger) bool {
	ok, err := contextAdmin(r)
	if !ok {
		if err != nil {
			http.Error(w, "User context not found", http.StatusUnauthorized)

			return false
		}

		log.Info("Not enough rights")

		http.Error(w, "Not enough rights", http.StatusForbidden)

		return false
	}
	return true
}

func changeField(w http.ResponseWriter, r *http.Request, log *slog.Logger, User AdminHandler, op, field string, value bool) {
	log.With(util.SlogWith(op, r)...)

	if !AdmCheck(w, r, log) {
		return
	}

	id := util.GetUrlParam(w, r, log)
	if id == 0 {
		log.Info("missing or wrong id")
		http.Error(w, "Missing or wrong id", http.StatusBadRequest)
		return
	}

	if n, err := User.UpdateField(field, id, value); err != nil {
		if n == 0 {
			log.Info(err.Error())

			http.Error(w, "No such user", http.StatusNotFound)
		} else if n == -2 {
			log.Info(err.Error())

			http.Error(w, "No such field", http.StatusBadRequest)
		}
		util.InternalError(w, r, log, err)
		return
	}

	user, err := User.Get(id)
	if err != nil {
		util.InternalError(w, r, log, err)
		return
	}

	log.Info(fmt.Sprintf("Successfully updated field: %v to %v", field, value))
	log.Debug(fmt.Sprintf("user: %v", user))

	render.JSON(w, r, user)
}
