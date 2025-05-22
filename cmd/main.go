package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"turbogate/config"
	"turbogate/internal/observability"
	"turbogate/internal/router"
	"turbogate/internal/watcher"
)

// Commented to commit the change for testing if github actions works
func main() {
	cfgPath := "config/config.yaml"

	//   Initialize Prometheus metrics first
	observability.InitMetrics()

	//   Load config
	routes, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	//   Set up rate limiter
	rateLimiter := router.NewRateLimiterMiddleware(5, 1)

	//   Set up initial router
	baseRouter := router.SetupRouter(routes, rateLimiter)

	//   Wrap router with Prometheus metrics middleware
	handlerWithMetrics := observability.MetricsMiddleware(baseRouter)

	//   Create ServeMux with both app routes and /metrics endpoint
	mux := http.NewServeMux()
	mux.Handle("/", handlerWithMetrics)
	mux.Handle("/metrics", observability.MetricsHandler())

	//   Create router manager for live reload
	manager := router.NewRouterManager(mux)

	//   Watch config file for changes
	go func() {
		err := watcher.WatchConfig(cfgPath, func() {
			newRoutes, err := config.LoadConfig(cfgPath)
			if err != nil {
				log.Printf("Error reloading config: %v", err)
				return
			}
			log.Println("  Reloaded routes")

			newRouter := router.SetupRouter(newRoutes, rateLimiter)
			handlerWithMetrics := observability.MetricsMiddleware(newRouter)

			newMux := http.NewServeMux()
			newMux.Handle("/", handlerWithMetrics)
			newMux.Handle("/metrics", observability.MetricsHandler())

			manager.UpdateHandler(newMux)
		})
		if err != nil {
			log.Fatalf("Watcher failed: %v", err)
		}
	}()

	//   Start HTTP server
	log.Println("🚀 TurboGate running at : 10000")
	server := &http.Server{
		Addr:    ":10000",
		Handler: manager,
	}

	//   Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("🛑 Shutting down...")
		server.Close()
	}()

	log.Fatal(server.ListenAndServe())
}
