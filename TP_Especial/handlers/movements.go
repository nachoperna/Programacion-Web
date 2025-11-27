package handlers

import (
	model "TP_Especial/structs"
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
	var movimientos []model.Movimientos
	if historial.LastDeposit.Valid {
		mov := model.Movimientos{Type: "Deposito", Amount: historial.LastDepositAmount.String, Date: historial.LastDeposit.Time}
		movimientos = append(movimientos, mov)
	}
	if historial.LastTransfer.Valid {
		mov := model.Movimientos{Type: "Transferencia", Amount: historial.LastTransferAmount.String, Date: historial.LastDeposit.Time}
		movimientos = append(movimientos, mov)
	}
	if historial.LastWithdrawal.Valid {
		mov := model.Movimientos{Type: "Retiro", Amount: historial.LastWithdrawalAmount.String, Date: historial.LastDeposit.Time}
		movimientos = append(movimientos, mov)
	}
	views.GetMovements(movimientos, alias).Render(h.ctx, w)
}
