package admin

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	r "github.com/sabbatD/srest-api/internal/lib/api/response"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
)

type AdminRequest u.Login

type AdminResponse struct {
	r.Response
}

type GetAllResponse struct {
	r.Response
	Users u.TableUser `json:"users,omitempty"`
}

type AdminHandler interface {
	UpdateField(field string, u u.Login, val any) error
	GetAll() ([]u.TableUser, error)
	Remove(u u.Login) error
}

func Block(log *slog.Logger, userHandler AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.hanlders.admin.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req AdminRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			return
		}

		// TODO: make block func call
	}
}

// TODO: Unblock http func

// TODO: MakeAdmin http func

// TODO: MakeUser http func

// TODO: Remove http func

// TODO: GetAll http func
