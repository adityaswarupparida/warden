package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 5
	config.MinConns = 2
	// config.MaxConnIdleTime = 10 * time.Second
	// config.MaxConnLifetime = 30 * time.Second
	return pgxpool.NewWithConfig(ctx, config)
}
