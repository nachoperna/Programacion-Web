package handlers

import (
	db "TP_Especial/db/sqlc"
	"TP_Especial/views"
	"encoding/csv"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
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

func (h *Handler) MovementsToCsv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=movements_report.csv")

	writer := csv.NewWriter(w)

	defer writer.Flush()

	writer.Write([]string{"TIPO", "MONTO", "DIA", "HORA"})
	historial, err := h.queries.GetHistoryComplete(h.ctx, r.URL.Query().Get("alias"))

	if err != nil {
		fmt.Printf("Error obteniendo historial de movimientos: %v", err)
	}

	for _, mov := range historial {
		writer.Write([]string{mov.Type, mov.Amount, mov.Day, mov.Time})
	}

}
