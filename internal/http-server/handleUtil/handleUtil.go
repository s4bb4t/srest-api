package handleutil

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Shortcut for logging
func SlogWith(op string, r *http.Request) []any {
	return []any{
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
		slog.String("\n", ""),
	}
}

// Shortcut for InternalError
func InternalError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.Debug(err.Error())

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// Shortcut for GetUrlParam
func GetUrlParam(w http.ResponseWriter, r *http.Request, log *slog.Logger) int {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		InternalError(w, r, log, err)
		return 0
	}
	return id
}
