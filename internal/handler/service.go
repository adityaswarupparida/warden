package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type CreateServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) GetServices(c *echo.Context) error {
	services, err := h.store.GetServices(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get services"})
	}

	return c.JSON(http.StatusOK, services)
}

func (h *Handler) CreateService(c *echo.Context) error {
	var req CreateServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.Name == "" || req.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name and description are required"})
	}

	svc, err := h.store.CreateService(c.Request().Context(), req.Name, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create service"})
	}

	return c.JSON(http.StatusCreated, svc)
}
