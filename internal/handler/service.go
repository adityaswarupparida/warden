package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

func (h *Handler) GetServices(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (h *Handler) CreateService(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (h *Handler) GetServiceIncidents(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (h *Handler) CreateServiceIncident(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
