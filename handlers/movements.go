package handlers

import (
	db "TP_Especial/db/sqlc"
	"TP_Especial/views"
	_ "github.com/lib/pq"
	"net/http"
)

func (h *Handler) ListMovements(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("alias")
	offset := h.GetOffset(r)
	historial, err := h.queries.GetHistory(h.ctx, db.GetHistoryParams{
		Alias:  alias,
		Offset: int32(offset),
	})
	if err != nil {
		http.Error(w, "Error obteniendo historial", http.StatusNotFound)
		return
	}
	siguientes, _ := h.queries.GetHistorySiguientes(h.ctx, db.GetHistorySiguientesParams{
		Alias:  alias,
		Offset: int32(offset + 3),
	})
	if siguientes == 0 {
		views.GetMovements(historial, alias, offset, false).Render(h.ctx, w)
	} else {
		views.GetMovements(historial, alias, offset, true).Render(h.ctx, w)
	}
}
