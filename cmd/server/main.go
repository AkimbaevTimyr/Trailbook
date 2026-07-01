package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"tracking-backend/internal/config"
	"tracking-backend/internal/db"
	httpDelivery "tracking-backend/internal/delivery/http"
	"tracking-backend/internal/delivery/http/handler"
	"tracking-backend/internal/repository"
	"tracking-backend/internal/server"
	"tracking-backend/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found or failed to load: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database.DB)
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUC)

	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authUC)
	authMW := httpDelivery.AuthMiddleware(authUC)

	r := httpDelivery.NewRouter(userHandler, authHandler, authMW)
	srv := server.New(cfg.Port, r)

	go func() {
		fmt.Printf("Server is running on port %s\n", cfg.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	fmt.Println("Server stopped")
}
