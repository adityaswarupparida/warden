package handler

import (
	"net/http"

	"github.com/adityaswarupparida/warden/internal/metrics"
	"github.com/labstack/echo/v5"
)

type CreateIncidentRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

func (h *Handler) GetServiceIncidents(c *echo.Context) error {
	serviceID := c.Param("id")

	incidents, err := h.store.GetIncidents(c.Request().Context(), serviceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get incidents"})
	}

	return c.JSON(http.StatusOK, incidents)
}

func (h *Handler) CreateServiceIncident(c *echo.Context) error {
	serviceID := c.Param("id")

	var req CreateIncidentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	inc, err := h.store.CreateIncident(c.Request().Context(), serviceID, req.Title, req.Description, req.Severity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create incident"})
	}

	metrics.ActiveIncidents.Inc()

	return c.JSON(http.StatusCreated, inc)
}

func (h *Handler) ResolveServiceIncident(c *echo.Context) error {
	serviceID := c.Param("id")
	incidentID := c.Param("incident_id")

	if err := h.store.ResolveIncident(c.Request().Context(), serviceID, incidentID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to resolve incident"})
	}

	metrics.ActiveIncidents.Dec()

	return c.NoContent(http.StatusNoContent)
}
