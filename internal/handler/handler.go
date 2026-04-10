package handler

import "github.com/adityaswarupparida/warden/internal/store"

type Handler struct {
	store store.Store
}

func NewHandler(store store.Store) *Handler {
	return &Handler{store: store}
}
