package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SKIND0A/online-shop/internal/config"
	"github.com/SKIND0A/online-shop/internal/delivery/http/handlers"
	"github.com/SKIND0A/online-shop/internal/repository/postgres"
	"github.com/SKIND0A/online-shop/internal/usecase"
)

func main() {
	cfg := config.Load()
	pool, err := postgres.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	//начало
	jwtTTL, err := time.ParseDuration(cfg.JWTAccessTTL)
	if err != nil {
		log.Fatalf("invalid JWT_ACCESS_TTL: %v", err)
	}

	userRepo := postgres.NewUserRepository(pool)
	jwtService := usecase.NewJWTService(cfg.JWTSecret, jwtTTL)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService)
	authHandler := handlers.NewAuthHandler(authUsecase)
	uiHandler := handlers.NewUIHandler()
	//конец

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, `{"success":false,"error":"db unavailable"}`, http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"status":"ok","db":"ok"},"error":null}`))
	})
	//начало
	http.HandleFunc("/api/v1/auth/register", authHandler.Register)
	http.HandleFunc("/api/v1/auth/login", authHandler.Login)
	http.HandleFunc("/", uiHandler.AuthPage)
	//конец
	fmt.Printf("server started on %s\n", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, nil); err != nil {
		log.Fatal(err)
	}
}
