package instrumentation

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

// PrometheusMiddleware is a middleware handler that records Prometheus metrics for each request.
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		statusCode := ww.Status()
		method := r.Method
		endpoint := r.URL.Path

		// Update Prometheus metrics
		if statusCode == 404 {
			ErrorCounter.WithLabelValues(method, endpoint).Inc()
		}
		RequestCounter.WithLabelValues(method).Inc()
		RequestDurationHistogram.WithLabelValues(method, endpoint).Observe(duration)

	})
}
