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
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var inc Incident
	err = tx.QueryRow(ctx,
		"INSERT INTO incidents (service_id, title, description, severity) VALUES ($1, $2, $3, $4) RETURNING id, service_id, title, description, status, severity, created_at, updated_at, resolved_at",
		serviceID, title, description, severity,
	).Scan(&inc.ID, &inc.ServiceID, &inc.Title, &inc.Description, &inc.Status, &inc.Severity, &inc.CreatedAt, &inc.UpdatedAt, &inc.ResolvedAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE services SET status = CASE
			WHEN EXISTS (SELECT 1 FROM incidents WHERE service_id = $1 AND resolved_at IS NULL AND severity = 'critical') THEN 'down'
			WHEN EXISTS (SELECT 1 FROM incidents WHERE service_id = $1 AND resolved_at IS NULL) THEN 'degraded'
			ELSE 'healthy'
		END
		WHERE id = $1`, serviceID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &inc, nil
}

func (s *pgStore) ResolveIncident(ctx context.Context, serviceID, incidentID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		"UPDATE incidents SET resolved_at = NOW(), updated_at = NOW() WHERE id = $1 AND service_id = $2 AND resolved_at IS NULL",
		incidentID, serviceID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE services SET status = CASE
			WHEN EXISTS (SELECT 1 FROM incidents WHERE service_id = $1 AND resolved_at IS NULL AND severity = 'critical') THEN 'down'
			WHEN EXISTS (SELECT 1 FROM incidents WHERE service_id = $1 AND resolved_at IS NULL) THEN 'degraded'
			ELSE 'healthy'
		END
		WHERE id = $1`, serviceID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
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
