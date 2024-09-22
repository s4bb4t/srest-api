package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	route := chi.NewRouter()
	route.Use(CORSMiddleware)

	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	route.Get("/static/*", http.HandlerFunc(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP))

	srv := &http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: route,
	}

	fmt.Println("Starting server on :8082")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("failed to start server:", err)
	}
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
