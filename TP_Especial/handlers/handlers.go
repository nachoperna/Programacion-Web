package handlers

import (
	sqlc "TP_Especial/db/sqlc"
	"context"
)

type Handler struct {
	queries *sqlc.Queries
	ctx     context.Context
}

func NewHandler(queries *sqlc.Queries, ctx context.Context) *Handler {
	return &Handler{
		queries: queries,
		ctx:     ctx,
	}
}
