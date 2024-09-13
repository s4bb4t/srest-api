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

type Users struct {
	Users []u.TableUser `json:"users"`
}

type AdminHandler interface {
	UpdateField(field string, id int, val any) (int64, error)
	GetAll(search, sotrOrder string, isBlocked bool, limit, offset int) ([]u.TableUser, error)
	Remove(id int) (int64, error)
	Get(id int) (u.TableUser, error)
	UpdateUser(u u.User, id int) (int64, error)
}

// Get godoc
// @Summary Get all users
// @Description Gets users by accepting a url query payload containing filters.
// Requires an Authorization header with a "Bearer token" for authentication.
// @Tags admin
// @Produce json
// @Param search query string false "Search in login or username or email"
// @Param sotrOrder query string false "sotrOrder asc or desc or none (asc, decs - sotrOrder by email, none - sotrOrder by id)"
// @Param isBlocked query bool false "block status"
// @Param limit query int false "limit of users for query"
// @Param offset query int false "offset"
// @Security BearerAuth
// @Success 200 {object} Users "Retrieve successful. Returns users."
// @Failure 401 {object} string "User context not found."
// @Failure 403 {object} string "Not enough rights."
// @Failure 500 {object} string "Internal error."
// @Router /admin/users [get]
func GetAll(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.GetAll"

		log.With(util.SlogWith(op, r)...)

		if !AdmCheck(w, r, log) {
			return
		}

		search := r.URL.Query().Get("search")
		sotrOrder := strings.ToLower(r.URL.Query().Get("sotrOrder"))
		isblockedStr := r.URL.Query().Get("isBlocked")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		switch sotrOrder {
		case "asc", "desc", "none":
		default:
			sotrOrder = "none"
		}
		if isblockedStr == "" {
			isblockedStr = "false"
		}
		if limitStr == "" {
			limitStr = "20"
		}
		if offsetStr == "" {
			offsetStr = "0"
		}

		isBlocked, _ := strconv.ParseBool(isblockedStr)
		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		users, err := User.GetAll(search, sotrOrder, isBlocked, limit, offset)
		if err != nil {
			util.InternalError(w, r, log, err)
			return
		}

		log.Info("users successfully retrieved")
		log.Debug(fmt.Sprintf("params: %v, %v, %v, %v, %v", search, sotrOrder, isBlocked, limit, offset))

		render.JSON(w, r, Users{users})
	}
}

// Profile godoc
// @Summary Retrieve user's profile
// @Description Retrieves user's profile by id.
// Requires an Authorization header with a "Bearer token" for authentication.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} u.TableUser "Retrieve successful. Returns user."
// @Failure 400 {object} string "Missing or wrong id."
// @Failure 401 {object} string "User context not found."
// @Failure 403 {object} string "Not enough rights."
// @Failure 404 {object} string "No such user."
// @Failure 500 {object} string "Internal error."
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
// @Summary Update user's fields
// @Description Updates user by id by accepting a JSON payload containing user.
// Requires an Authorization header with a "Bearer token" for authentication.
// @Tags admin
// @Accept json
// @Produce json
// @Param UserData body u.User true "Complete user data"
// @Security BearerAuth
// @Success 200 {object} u.TableUser "Update successful. Returns user ok."
// @Failure 400 {object} string "failed to deserialize json request."
// @Failure 400 {object} string "Missing or wrong id."
// @Failure 400 {object} string "Login or email already used."
// @Failure 401 {object} string "User context not found."
// @Failure 403 {object} string "Not enough rights."
// @Failure 404 {object} string "No such user."
// @Failure 500 {object} string "Internal error."
// @Router /admin/users/{id} [put]
func UpdateUser(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.UpdateUser"

		log.With(util.SlogWith(op, r)...)

		if !AdmCheck(w, r, log) {
			return
		}

		var req u.User
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
			} else if n == -2 {

				log.Info(err.Error())

				http.Error(w, "Login or email already used", http.StatusBadRequest)
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
// @Description Removes user by id in url.
// Requires an Authorization header with a "Bearer token" for authentication.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} string "Remove successful. Returns ok."
// @Failure 400 {object} string "Missing or wrong id."
// @Failure 401 {object} string "User context not found."
// @Failure 403 {object} string "Not enough rights."
// @Failure 404 {object} string "No such user."
// @Failure 500 {object} string "Internal error."
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
// @Description Blocks user by id in url.
// Requires an Authorization header with a "Bearer token" for authentication.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} u.TableUser "Block successful. Returns user ok."
// @Failure 400 {object} string "Missing or wrong id."
// @Failure 400 {object} string "No such field."
// @Failure 404 {object} string "No such user."
// @Failure 500 {object} string "Internal error."
// @Router /admin/users/{id}/block [post]
func Block(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.Block"

		changeField(w, r, log, User, op, "block", true)
	}
}

// Unlock godoc
// @Summary Unlock user
// @Description Unlocks user by id in url.
// Requires an Authorization header with a "Bearer token" for authentication.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} u.TableUser "Unlock successful. Returns user ok."
// @Failure 400 {object} string "Missing or wrong id."
// @Failure 400 {object} string "No such field."
// @Failure 404 {object} string "No such user."
// @Failure 500 {object} string "Internal error."
// @Router /admin/users/{id}/unlock [post]
func Unblock(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.Unblock"

		changeField(w, r, log, User, op, "block", false)
	}
}

// Update godoc
// @Summary Update user's rights
// @Description Updates user by id by accepting a JSON payload containing user's rights.
// @Tags admin
// @Accept json
// @Produce json
// @Param UserData body UpdateRequest true "Complete user data"
// @Success 200 {object} u.TableUser "Update successful. Returns ok."
// @Failure 400 {object} string "failed to deserialize json request."
// @Failure 400 {object} string "Missing or wrong id."
// @Failure 400 {object} string "No such field."
// @Failure 404 {object} string "No such user."
// @Failure 500 {object} string "Internal error."
// @Router /admin/users/{id}/rights [post]
func Update(log *slog.Logger, User AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.admin.Update"

		log.With(util.SlogWith(op, r)...)

		// if !AdmCheck(w, r, log) {
		// 	return
		// }

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

		log.Info("Successfully updated user's rights")
		log.Debug(fmt.Sprintf("user: %v", user))

		render.JSON(w, r, user)
	}
}

func contextAdmin(r *http.Request) (bool, error) {
	userContext, ok := r.Context().Value(access.CxtKey("userContext")).(access.UserContext)
	if !ok {
		return false, fmt.Errorf("Unauthorized")
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
