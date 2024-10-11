package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func main() {
	route := chi.NewRouter()
	route.Use(CORSMiddleware)

	// Отдаем только index.html
	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("/usr/local/bin/frontend/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Обработчик для доступа к статическим файлам только для самого сервиса
	route.Handle("/*", http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isInternalRequest(r) {
			http.FileServer(http.Dir("/usr/local/bin/frontend")).ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})))

	srv := &http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: route,
	}

	fmt.Println("Starting server on :8081")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("failed to start server:", err)
	}
}

// Функция для проверки, является ли запрос локальным
func isInternalRequest(r *http.Request) bool {
	ip := r.RemoteAddr
	parts := strings.Split(ip, ":")
	if len(parts) > 0 {
		if parts[0] == "127.0.0.1" || parts[0] == "::1" {
			return true // Запрос с локального адреса
		}
	}
	return false // Запрос с внешнего адреса
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
