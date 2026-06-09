// Package metrics declares the Prometheus registry and the HTTP instruments
// (counter + histogram) shared across the service, plus the recording middleware.
package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Registry is the single source of truth for every metric in this service.
var Registry = prometheus.NewRegistry()

// httpRequestsTotal is a monotonically increasing counter consumed via rate().
var httpRequestsTotal = promauto.With(Registry).NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests",
	},
	[]string{"method", "route", "status_code"},
)

// httpRequestDurationSeconds is a histogram whose fixed buckets let PromQL
// compute quantiles (p95, p99) at query time via histogram_quantile.
var httpRequestDurationSeconds = promauto.With(Registry).NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
	},
	[]string{"method", "route"},
)

func init() {
	// Register Go runtime + process gauges (goroutines, heap, GC) on the same registry.
	Registry.MustRegister(collectors.NewGoCollector())
	Registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}

// statusRecorder wraps ResponseWriter to capture the final status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Middleware records the counter + histogram after the wrapped handler finishes.
// The route argument is the bounded route template, never the raw path.
func Middleware(next http.Handler, route string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		// Defer runs after the handler finishes, so we read the real status.
		defer func() {
			durationSec := time.Since(start).Seconds()
			status := strconv.Itoa(rec.status)
			httpRequestsTotal.WithLabelValues(req.Method, route, status).Inc()
			httpRequestDurationSeconds.WithLabelValues(req.Method, route).Observe(durationSec)
		}()
		next.ServeHTTP(rec, req)
	})
}
