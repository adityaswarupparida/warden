package handler

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/adityaswarupparida/warden/internal/store"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

// GetServiceIncidents
func TestGetServiceIncidents_Returns200(t *testing.T) {
	h := NewHandler(&mockStore{incidents: []store.Incident{}})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodGet, "/api/v1/services/svc-1/incidents", "")
	c.Set("id", "svc-1")

	_ = h.GetServiceIncidents(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetServiceIncidents_StoreFailure_Returns500(t *testing.T) {
	h := NewHandler(&mockStore{err: errors.New("db error")})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodGet, "/api/v1/services/svc-1/incidents", "")
	c.Set("id", "svc-1")

	_ = h.GetServiceIncidents(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to get incidents")
}

// CreateServiceIncident
func TestCreateServiceIncident_Returns201(t *testing.T) {
	h := NewHandler(&mockStore{
		incident: &store.Incident{ID: "1", ServiceID: "svc-1", Title: "DB down", Status: "identified", Severity: "high", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
	})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodPost, "/api/v1/services/svc-1/incidents", `{"title":"DB down","description":"unreachable","severity":"high"}`)
	c.Set("id", "svc-1")

	_ = h.CreateServiceIncident(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateServiceIncident_InvalidBody_Returns400(t *testing.T) {
	h := NewHandler(&mockStore{})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodPost, "/api/v1/services/svc-1/incidents", `invalid json`)
	c.Set("id", "svc-1")

	_ = h.CreateServiceIncident(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid request")
}

func TestCreateServiceIncident_StoreFailure_Returns500(t *testing.T) {
	h := NewHandler(&mockStore{err: errors.New("db error")})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodPost, "/api/v1services/svc-1/incidents", `{"title":"DB down","description":"unreachable","severity":"high"}`)
	c.Set("id", "svc-1")

	_ = h.CreateServiceIncident(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to create incident")
}
