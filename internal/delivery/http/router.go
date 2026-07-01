package http

import (
	"net/http"

	"tracking-backend/internal/delivery/http/handler"
	"tracking-backend/internal/delivery/http/middleware"
	"tracking-backend/internal/delivery/http/response"
)

func NewRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, authMW func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, map[string]string{"status": "ok"})
	})

	mux.Handle("GET /health/auth", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, map[string]string{"status": "authenticated"})
	})))

	mux.HandleFunc("POST /login", authHandler.Login)

	mux.HandleFunc("GET /users", userHandler.GetUsers)
	mux.Handle("GET /users/{id}", authMW(http.HandlerFunc(userHandler.GetUser)))
	mux.HandleFunc("POST /users", userHandler.CreateUser)

	return mux
}

var AuthMiddleware = middleware.Auth
