package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"turbogate/config"
)

var routeMap map[string]string
var ReloadRouter func()

type AddRouteRequest struct {
	Route  string `json:"route"`
	Domain string `json:"domain"`
}

func AddRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("route")
	target := r.URL.Query().Get("domain")

	if path == "" || target == "" {
		http.Error(w, "Missing route or domain", http.StatusBadRequest)
		return
	}

	cfgPath := "config/config.yaml"
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	// Check for duplicates
	for _, r := range cfg.Routes {
		if r.Path == path {
			http.Error(w, "Route already exists", http.StatusConflict)
			return
		}
	}

	// Append in the yaml and save
	cfg.Routes = append(cfg.Routes, config.Route{Path: path, Target: target})
	err = config.SaveConfig(cfgPath, cfg)
	if err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	// 🔁 Call the hot-reload function from main
	if ReloadRouter != nil {
		ReloadRouter()
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Route added and router reloaded"))

}

func ReloadRoutes(cfgPath string) error {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	newRouteMap := make(map[string]string)
	for _, r := range cfg.Routes {
		newRouteMap[r.Path] = r.Target
	}

	// Replace global route map atomically
	routeMap = newRouteMap

	return nil
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	target, ok := routeMap[r.URL.Path]
	if !ok {
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}

	remote, err := url.Parse(target)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}
