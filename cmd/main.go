package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"TurboGate/config"
	"TurboGate/internal/observability"
	"TurboGate/internal/router"
	"TurboGate/internal/watcher"
)

func main() {
	cfgPath := "config/config.yaml"
	routes, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	rateLimiter := router.NewRateLimiterMiddleware(10, 20)

	// Set up initial router
	baseRouter := router.SetupRouter(routes, rateLimiter)

	// Wrap it with metrics
	handlerWithMetrics := observability.MetricsMiddleware(baseRouter)

	// Set up metrics endpoint separately
	mux := http.NewServeMux()
	mux.Handle("/", handlerWithMetrics)
	mux.Handle("/metrics", observability.MetricsHandler())

	// Graceful hot reload setup
	manager := router.NewRouterManager(mux)

	go func() {
		err := watcher.WatchConfig(cfgPath, func() {
			newRoutes, err := config.LoadConfig(cfgPath)
			if err != nil {
				log.Printf("Error reloading config: %v", err)
				return
			}
			log.Println("✅ Reloaded routes")

			newRouter := router.SetupRouter(newRoutes, rateLimiter)
			handlerWithMetrics := observability.MetricsMiddleware(newRouter)

			manager.UpdateHandler(handlerWithMetrics)
		})
		if err != nil {
			log.Fatalf("Watcher failed: %v", err)
		}
	}()

	observability.InitMetrics()

	log.Println("🚀 TurboGate running at :8080")
	server := &http.Server{
		Addr:    ":8080",
		Handler: manager,
	}

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("🛑 Shutting down...")
		server.Close()
	}()

	log.Fatal(server.ListenAndServe())
}
