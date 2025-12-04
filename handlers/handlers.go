package handlers

import (
	sqlc "TP_Especial/db/sqlc"
	"context"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Handler struct {
	queries       *sqlc.Queries
	ctx           context.Context
	RowLimitTable int
}

func NewHandler(queries *sqlc.Queries, ctx context.Context) *Handler {
	return &Handler{
		queries:       queries,
		ctx:           ctx,
		RowLimitTable: 3,
	}
}

func (h *Handler) GetOffset(r *http.Request) int {
	orden := r.URL.Query().Get("order")
	offset := r.URL.Query().Get("offset")
	aux := 0
	if orden != "" && offset != "" {
		off, _ := strconv.Atoi(offset)
		if orden == "anterior" {
			aux = int(off - h.RowLimitTable)
		} else {
			aux = int(off + h.RowLimitTable)
		}
	}
	return aux
}

func (h *Handler) balanceFormateado(balance string) string {
	printer := message.NewPrinter(language.Spanish)
	balance_float, _ := strconv.ParseFloat(balance, 64)
	balance_formateado := printer.Sprintf("%.2f", balance_float)
	return balance_formateado
}
