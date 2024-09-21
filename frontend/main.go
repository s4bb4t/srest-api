package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	route := chi.NewRouter()
	route.Route("/", func(router chi.Router) {
		router.Use(CORSMiddleware)

		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			t, err := template.ParseFiles("index.html")
			if err != nil {
				fmt.Fprint(w, err.Error())
			}

			t.ExecuteTemplate(w, "main", nil)
		})

		srv := &http.Server{
			Addr:    "0.0.0.0:8082",
			Handler: route,
		}

		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("failed to start server")
		}

		fmt.Println("server stopped")
	})

}
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
