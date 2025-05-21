package main

import (
	"log"
	"net/http"

	"TurboGate/config"
	"TurboGate/internal/limiter"
	"TurboGate/internal/router"
)

func main() {
	// Load routing config from YAML
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create per-IP rate limiter: 100 req/min
	lim := limiter.NewIPRateLimiter(100, 60)

	// Create router with config and rate limiter
	mux := router.NewRouter(cfg, lim)

	log.Println("TurboGate listening on :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
