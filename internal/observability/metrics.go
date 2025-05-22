package observability

import (
	"net/http"
	"strconv"
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

		duration := time.Since(start).Seconds()
		requestCount.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(rw.status)).Inc()
		requestDuration.WithLabelValues(r.URL.Path).Observe(duration)
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
