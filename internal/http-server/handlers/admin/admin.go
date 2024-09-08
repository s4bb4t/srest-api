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

// Register godoc
// @Summary Update user's rights
// @Description Updates user by id by accepting a JSON payload containing user's rights.
// @Tags admin
// @Accept json
// @Produce json
// @Param UserData body UpdateRequest true "Complete user data"
// @Success 200 {object} resp.Response "Update successful. Returns user ok."
// @Failure 401 {object} resp.Response "Update failed. Returns error message."
// @Router /admin/users/{id}/rights [post]
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

// Register godoc
// @Summary Block user
// @Description Blocks user by id in url.
// @Tags admin
// @Produce json
// @Success 200 {object} resp.Response "Block successful. Returns user ok."
// @Failure 401 {object} resp.Response "Block failed. Returns error message."
// @Router /admin/users/{id}/block [post]
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

		log.Info(fmt.Sprintf("Successfully updated field: %v to %v", "block", true))

		render.JSON(w, r, resp.OK())
	}
}

// Register godoc
// @Summary Unlock user
// @Description Unlocks user by id in url.
// @Tags admin
// @Produce json
// @Success 200 {object} resp.Response "Unlock successful. Returns user ok."
// @Failure 401 {object} resp.Response "Unlock failed. Returns error message."
// @Router /admin/users/{id}/unlock [post]
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

		log.Info(fmt.Sprintf("Successfully updated field: %v to %v", "block", false))

		render.JSON(w, r, resp.OK())
	}
}

// Register godoc
// @Summary Remove user
// @Description Removes user by id in url.
// @Tags admin
// @Produce json
// @Success 200 {object} resp.Response "Remove successful. Returns user ok."
// @Failure 401 {object} resp.Response "Remove failed. Returns error message."
// @Router /admin/users/{id} [delete]
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

// Register godoc
// @Summary Get user's rights
// @Description Gets users by accepting a url query payload containing filters.
// @Tags admin
// @Produce json
// @Param search query string false "Search term"
// @Param order query string true "order asc or desc"
// @Param blocked query bool true "block status"
// @Param limit query int true "limit of users for query"
// @Param offset query int true "offset"
// @Success 200 {object} GetAllResponse "Retrieve successful. Returns users."
// @Failure 401 {object} resp.Response "Retrieve failed. Returns error message."
// @Router /admin/users [get]
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

		if order == "" || blockedStr == "" || limitStr == "" || offsetStr == "" {
			log.Info("one or more parameters are empty")

			render.JSON(w, r, resp.Error("one or more parameters are empty"))

			return
		}

		blocked, _ := strconv.ParseBool(blockedStr)
		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

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

// Register godoc
// @Summary Retrieve user's profile
// @Description Retrieves user's profile by id.
// @Tags admin
// @Produce json
// @Success 200 {object} GetResponse "Retrieve successful. Returns user."
// @Failure 401 {object} resp.Response "Retrieve failed. Returns error message."
// @Router /admin/users/{id} [get]
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

// Register godoc
// @Summary Update user's rights
// @Description Updates user by id by accepting a JSON payload containing user's rights.
// @Tags admin
// @Accept json
// @Produce json
// @Param UserData body u.User true "Complete user data"
// @Success 200 {object} resp.Response "Update successful. Returns user ok."
// @Failure 401 {object} resp.Response "Update failed. Returns error message."
// @Router /admin/users/{id}/rights [put]
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
