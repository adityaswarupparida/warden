package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Incident struct {
	ID          string     `db:"id"`
	ServiceID   string     `db:"service_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	Severity    string     `db:"severity"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	ResolvedAt  *time.Time `db:"resolved_at"`
}

type IncidentUpdate struct {
	ID         string    `db:"id"`
	IncidentID string    `db:"incident_id"`
	Message    string    `db:"message"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
}

func (s *pgStore) GetIncidents(ctx context.Context, serviceID string) ([]Incident, error) {
	rows, err := s.pool.Query(ctx,
		"SELECT id, service_id, title, description, status, severity, created_at, updated_at, resolved_at FROM incidents WHERE service_id = $1",
		serviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[Incident])
}

func (s *pgStore) CreateIncident(ctx context.Context, serviceID, title, description, severity string) (*Incident, error) {
	var inc Incident
	err := s.pool.QueryRow(ctx,
		"INSERT INTO incidents (service_id, title, description, severity) VALUES ($1, $2, $3, $4) RETURNING id, service_id, title, description, status, severity, created_at, updated_at, resolved_at",
		serviceID, title, description, severity,
	).Scan(&inc.ID, &inc.ServiceID, &inc.Title, &inc.Description, &inc.Status, &inc.Severity, &inc.CreatedAt, &inc.UpdatedAt, &inc.ResolvedAt)
	if err != nil {
		return nil, err
	}
	return &inc, nil
}

func (s *pgStore) CreateIncidentUpdate(ctx context.Context, incidentID, message, status string) (*IncidentUpdate, error) {
	var u IncidentUpdate
	err := s.pool.QueryRow(ctx,
		"INSERT INTO incident_updates (incident_id, message, status) VALUES ($1, $2, $3) RETURNING id, incident_id, message, status, created_at",
		incidentID, message, status,
	).Scan(&u.ID, &u.IncidentID, &u.Message, &u.Status, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
