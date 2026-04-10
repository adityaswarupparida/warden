package api

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	errorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP errors",
	}, []string{"method", "path", "status"})
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Path() == "/metrics" || c.Path() == "/healthz" {
				return next(c)
			}

			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			_, status := echo.ResolveResponseStatus(c.Response(), err)

			statusStr := strconv.Itoa(status)
			requestsTotal.WithLabelValues(c.Request().Method, c.Path(), statusStr).Inc()
			if status >= 400 {
				errorsTotal.WithLabelValues(c.Request().Method, c.Path(), statusStr).Inc()
			}
			requestDuration.WithLabelValues(c.Request().Method, c.Path()).Observe(duration.Seconds())

			return err
		}
	}
}
