package http

import (
	"net/http"

	"tracking-backend/internal/delivery/http/handler"
	"tracking-backend/internal/delivery/http/middleware"
	"tracking-backend/internal/delivery/http/response"
)

// NewRouter builds and returns the application HTTP router.
func NewRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, routeHandler *handler.RouteHandler, authMW func(http.Handler) http.Handler) http.Handler {
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

	mux.HandleFunc("GET /routes", routeHandler.GetRoutes)
	mux.Handle("GET /routes/my", authMW(http.HandlerFunc(routeHandler.GetMyRoutes)))
	mux.HandleFunc("GET /routes/{id}", routeHandler.GetRoute)
	mux.Handle("POST /routes", authMW(http.HandlerFunc(routeHandler.CreateRoute)))
	mux.HandleFunc("PUT /routes/{id}", routeHandler.UpdateRoute)
	mux.HandleFunc("DELETE /routes/{id}", routeHandler.DeleteRoute)

	return mux
}

// AuthMiddleware is an alias for middleware.Auth to keep wiring concise.
var AuthMiddleware = middleware.Auth
