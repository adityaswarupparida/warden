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

// GetServices
func TestGetServices_Returns200(t *testing.T) {
	h := NewHandler(&mockStore{services: []store.Service{}})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodGet, "/api/v1/services", "")

	_ = h.GetServices(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetServices_StoreFailure_Returns500(t *testing.T) {
	h := NewHandler(&mockStore{err: errors.New("db connection lost")})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodGet, "/api/v1/services", "")

	_ = h.GetServices(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to get services")
}

// CreateService
func TestCreateService_Returns201(t *testing.T) {
	h := NewHandler(&mockStore{
		service: &store.Service{ID: "1", Name: "API", Status: "healthy", CreatedAt: time.Now().UTC()},
	})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodPost, "/api/v1/services", `{"name":"API","description":"main"}`)

	_ = h.CreateService(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateService_MissingBody_Returns400(t *testing.T) {
	h := NewHandler(&mockStore{})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodPost, "/api/v1/services", `invalid json`)

	_ = h.CreateService(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid request")
}

func TestCreateService_BlankValues_Returns400(t *testing.T) {
	h := NewHandler(&mockStore{})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodPost, "/api/v1/services", `{"name":"","description":""}`)

	_ = h.CreateService(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "name and description are required")
}

func TestCreateService_StoreFailure_Returns500(t *testing.T) {
	h := NewHandler(&mockStore{err: errors.New("db error")})
	e := echo.New()
	c, rec := newTestContext(e, http.MethodPost, "/api/v1/services", `{"name":"API","description":"main"}`)

	_ = h.CreateService(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to create service")
}
