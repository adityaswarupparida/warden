package api

import (
	"net/http"

	"github.com/labstack/echo/v5"
	// "github.com/labstack/echo/v5/middleware"
	"github.com/adityaswarupparida/warden/internal/handler"
)

func RegisterRoutes(e *echo.Echo, h *handler.Handler) {
	// e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	apiV1 := e.Group("/api/v1")
	// Services routes
	apiV1.GET("/services", h.GetServices)
	apiV1.POST("/services", h.CreateService)
	apiV1.GET("/services/:id/incidents", h.GetServiceIncidents)
	apiV1.POST("/services/:id/incidents", h.CreateServiceIncident)

}

// func ListenAndServe(e *echo.Echo) error {
// 	if err := e.Start(":1323"); err != nil {
// 		e.Logger.Error("failed to start server", "error", err)
// 		return err
// 	}
// 	return nil
// }
