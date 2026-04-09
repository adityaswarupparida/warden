package main

import (
	"context"
	"log"
	"os"

	"github.com/adityaswarupparida/warden/internal/api"
	"github.com/adityaswarupparida/warden/internal/handler"
	"github.com/adityaswarupparida/warden/internal/store"

	"github.com/labstack/echo/v5"
)

func main() {
	ctx := context.Background()
	pool, err := store.NewPool(ctx, os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	h := handler.NewHandler(pool)

	e := echo.New()
	api.RegisterRoutes(e, h)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
