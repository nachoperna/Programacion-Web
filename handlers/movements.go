package handlers

import (
	db "TP_Especial/db/sqlc"
	"TP_Especial/views"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func (h *Handler) ListMovements(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("alias")
	orden := r.URL.Query().Get("order")
	offset := r.URL.Query().Get("offset")
	aux := 0
	limit := 3
	if orden != "" && offset != "" {
		off, _ := strconv.Atoi(offset)
		if orden == "anterior" {
			aux = int(off - limit)
		} else {
			aux = int(off + limit)
		}
	}
	historial, err := h.queries.GetHistory(h.ctx, db.GetHistoryParams{
		Alias:  alias,
		Offset: int32(aux),
	})
	if err != nil {
		http.Error(w, "Error obteniendo historial", http.StatusNotFound)
		return
	}
	siguientes, _ := h.queries.GetHistorySiguientes(h.ctx, db.GetHistorySiguientesParams{
		Alias:  alias,
		Offset: int32(aux + 3),
	})
	if siguientes == 0 {
		views.GetMovements(historial, alias, aux, false).Render(h.ctx, w)
	} else {
		views.GetMovements(historial, alias, aux, true).Render(h.ctx, w)
	}
}
