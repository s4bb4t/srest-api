package admin

import (
	"log/slog"
	"net/http"

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
	GetAllUsers() ([]u.TableUser, error)
}

func New(log *slog.Logger, userHandler AdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
