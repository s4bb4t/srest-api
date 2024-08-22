package user

import (
	"log/slog"
	"net/http"

	r "github.com/sabbatD/srest-api/internal/lib/api/response"
	u "github.com/sabbatD/srest-api/internal/lib/userConfig"
)

type RegisterRequest u.User

type RegisterResponse struct {
	r.Response
	Authdata u.AuthData `json:"authdata,omitempty"`
}

type AuthRequest u.AuthData

type AuthResponse struct {
	r.Response
}

type UserHandler interface {
	Auth(u u.AuthData) (bool, error)
}

func New(log *slog.Logger, userHandler UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
