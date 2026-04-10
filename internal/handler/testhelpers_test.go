package handler

import (
	"context"
	"net/http/httptest"
	"strings"

	"github.com/adityaswarupparida/warden/internal/store"
	"github.com/labstack/echo/v5"
)

type mockStore struct {
	services  []store.Service
	service   *store.Service
	incident  *store.Incident
	incidents []store.Incident
	err       error
}

func (m *mockStore) GetServices(ctx context.Context) ([]store.Service, error) {
	return m.services, m.err
}
func (m *mockStore) CreateService(ctx context.Context, name, description string) (*store.Service, error) {
	return m.service, m.err
}
func (m *mockStore) GetIncidents(ctx context.Context, serviceID string) ([]store.Incident, error) {
	return m.incidents, m.err
}
func (m *mockStore) CreateIncident(ctx context.Context, serviceID, title, description, severity string) (*store.Incident, error) {
	return m.incident, m.err
}

func newTestContext(e *echo.Echo, method, path, body string) (*echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
