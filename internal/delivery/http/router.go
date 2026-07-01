package http

import (
	"net/http"

	"tracking-backend/internal/delivery/http/handler"
	"tracking-backend/internal/delivery/http/middleware"
	"tracking-backend/internal/delivery/http/response"
)

// NewRouter builds and returns the application HTTP router.
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
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("POST /users", userHandler.CreateUser)

	return mux
}

// AuthMiddleware is an alias for middleware.Auth to keep wiring concise.
var AuthMiddleware = middleware.Auth
