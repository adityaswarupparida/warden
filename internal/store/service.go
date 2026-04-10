package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}

func (s *pgStore) GetServices(ctx context.Context) ([]Service, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name, description, status, created_at FROM services")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[Service])
}

func (s *pgStore) CreateService(ctx context.Context, name string, description string) (*Service, error) {
	var svc Service
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO services (name, description) VALUES ($1, $2) RETURNING id, name, description, status, created_at",
		name, description,
	).Scan(&svc.ID, &svc.Name, &svc.Description, &svc.Status, &svc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &svc, nil
}
