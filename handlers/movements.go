package handlers

import (
	"TP_Especial/views"
	_ "github.com/lib/pq"
	"net/http"
)

func (h *Handler) ListMovements(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("alias")
	historial, err := h.queries.GetHistory(h.ctx, alias)
	if err != nil {
		http.Error(w, "Error obteniendo requests", http.StatusNotFound)
		return
	}
	views.GetMovements(historial).Render(h.ctx, w)
}
