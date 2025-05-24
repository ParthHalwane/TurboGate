package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"turbogate/config"
	"turbogate/internal/api"
	"turbogate/internal/observability"
	"turbogate/internal/router"
	"turbogate/internal/watcher"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Initialize Prometheus metrics first
	observability.InitMetrics()

	// Set up rate limiter
	currentRateLimiter = router.NewRateLimiterMiddleware(100000, 100000)

	cfgPath := "config/config.yaml"
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Set up initial router
	baseRouter := router.SetupRouter(cfg.Routes, currentRateLimiter)

	// Wrap router with Prometheus metrics middleware
	handlerWithMetrics = observability.MetricsMiddleware(baseRouter)
	api.ReloadRouter = reloadRouter

	// Create ServeMux with both app routes and /metrics endpoint
	mux := http.NewServeMux()
	wwr := 1
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "application/json")
			log.Println("Backend running", wwr)
			wwr += 1
			json.NewEncoder(w).Encode(map[string]string{"status": "TurboGate backend running 🚀"})
			return
		}
		handlerWithMetrics.ServeHTTP(w, r)
	})
	mux.Handle("/metrics", observability.MetricsHandler())
	mux.HandleFunc("/api/add-route", api.AddRouteHandler)
	mux.HandleFunc("/api/routes", ListRoutesHandler)

	handlerWithCORS := corsMiddleware(mux)

	// Create router manager for live reload
	manager := router.NewRouterManager(handlerWithCORS)

	// Watch config file for changes
	go func() {
		err := watcher.WatchConfig(cfgPath, func() {
			newCfg, err := config.LoadConfig(cfgPath)
			if err != nil {
				log.Printf("❌ Error reloading config: %v", err)
				return
			}
			log.Println("🔁 Reloaded routes")

			newRouter := router.SetupRouter(newCfg.Routes, currentRateLimiter)
			newHandler := observability.MetricsMiddleware(newRouter)

			newMux := http.NewServeMux()
			newMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/" {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]string{"status": "TurboGate backend running 🚀"})
					return
				}
				newHandler.ServeHTTP(w, r) // Forward to your full app router
			})
			newMux.Handle("/metrics", observability.MetricsHandler())
			newMux.HandleFunc("/api/add-route", api.AddRouteHandler)
			newMux.HandleFunc("/api/routes", ListRoutesHandler)

			handlerWithCORS := corsMiddleware(newMux)
			manager.UpdateHandler(handlerWithCORS)
		})
		if err != nil {
			log.Fatalf("❌ Watcher failed: %v", err)
		}
	}()

	// Start HTTP server
	server := &http.Server{
		Addr:           ":10000",
		Handler:        manager,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("🚀 TurboGate running at :10000")

	http.DefaultTransport.(*http.Transport).MaxIdleConns = 10000
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 10000
	http.DefaultTransport.(*http.Transport).IdleConnTimeout = 90 * time.Second

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v", err)
	}

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("🛑 Shutting down...")
		server.Close()
	}()

	log.Fatal(server.ListenAndServe())
}

func ListRoutesHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		http.Error(w, "Failed to load routes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cfg.Routes)
	if err != nil {
		http.Error(w, "Failed to encode routes", http.StatusInternalServerError)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var (
	handlerWithMetrics http.Handler
	currentRateLimiter func(http.Handler) http.Handler
)

func reloadRouter() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Println("❌ Failed to reload config:", err)
		return
	}

	newRouter := router.SetupRouter(cfg.Routes, currentRateLimiter)
	handlerWithMetrics = observability.MetricsMiddleware(newRouter)
	log.Println("✅ Router reloaded successfully")
}
