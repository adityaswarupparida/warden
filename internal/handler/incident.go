package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type CreateIncidentRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// type CreateIncidentUpdateRequest struct {
// 	Message string `json:"message"`
// 	Status  string `json:"status"`
// }

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

	return c.JSON(http.StatusCreated, inc)
}

// func (h *Handler) CreateIncidentUpdate(c *echo.Context) error {
// 	incidentID := c.Param("incident_id")

// 	var req CreateIncidentUpdateRequest
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
// 	}

// 	update, err := h.store.CreateIncidentUpdate(c.Request().Context(), incidentID, req.Message, req.Status)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create incident update"})
// 	}

// 	return c.JSON(http.StatusCreated, update)
// }
