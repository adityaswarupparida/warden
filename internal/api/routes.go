package api

import (
	"net/http"

	"github.com/adityaswarupparida/warden/internal/handler"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRoutes(e *echo.Echo, h *handler.Handler) {
	e.Use(middleware.RequestLogger())
	e.Use(PrometheusMiddleware())

	e.GET("/healthz", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	apiV1 := e.Group("/api/v1")
	// Services routes
	apiV1.GET("/services", h.GetServices)
	apiV1.POST("/services", h.CreateService)
	apiV1.GET("/services/:id/incidents", h.GetServiceIncidents)
	apiV1.POST("/services/:id/incidents", h.CreateServiceIncident)

}
