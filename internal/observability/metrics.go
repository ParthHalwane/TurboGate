package observability

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "turbogate_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "turbogate_http_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	inFlightRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "turbogate_http_inflight_requests",
			Help: "Current number of in-flight requests",
		},
	)

	requestFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "turbogate_http_request_failures_total",
			Help: "Total failed HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(requestCount, requestDuration, inFlightRequests)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inFlightRequests.Inc()
		defer inFlightRequests.Dec()

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)
		path := getPathGroup(r.URL.Path)
		duration := time.Since(start).Seconds()

		if rw.status >= 500 {
			requestFailures.WithLabelValues(path, r.Method, strconv.Itoa(rw.status)).Inc()
		}

		requestCount.WithLabelValues(path, r.Method, strconv.Itoa(rw.status)).Inc()
		requestDuration.WithLabelValues(path).Observe(duration)
	})
}

// Custom responseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Expose Prometheus metrics endpoint
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func getPathGroup(path string) string {
	segments := strings.Split(path, "/")
	if len(segments) > 2 {
		segments[2] = ":id"
	}
	return strings.Join(segments, "/")
}
