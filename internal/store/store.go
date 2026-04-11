package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	GetServices(ctx context.Context) ([]Service, error)
	CreateService(ctx context.Context, name, description string) (*Service, error)
	GetIncidents(ctx context.Context, serviceID string) ([]Incident, error)
	CreateIncident(ctx context.Context, serviceID, title, description, severity string) (*Incident, error)
	ResolveIncident(ctx context.Context, serviceID, incidentID string) error
}

type pgStore struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &pgStore{pool: pool}
}
