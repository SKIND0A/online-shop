package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SKIND0A/online-shop/internal/config"
	"github.com/SKIND0A/online-shop/internal/repository/postgres"
)

func main() {
	cfg := config.Load()
	pool, err := postgres.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, `{"success":false,"error":"db unavailable"}`, http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"status":"ok","db":"ok"},"error":null}`))
	})
	fmt.Printf("server started on %s\n", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, nil); err != nil {
		log.Fatal(err)
	}
}
