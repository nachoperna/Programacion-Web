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
	if historial.LastDepositAmount.Valid {
		mov := model.Movimientos{Type: "Deposito", Amount: historial.LastDepositAmount.String, Day: historial.DepositDay, Time: historial.DepositTime}
		movimientos = append(movimientos, mov)
	}
	if historial.LastTransferAmount.Valid {
		mov := model.Movimientos{Type: "Transferencia", Amount: historial.LastTransferAmount.String, Day: historial.TransferDay, Time: historial.TransferTime}
		movimientos = append(movimientos, mov)
	}
	if historial.LastWithdrawalAmount.Valid {
		mov := model.Movimientos{Type: "Retiro", Amount: historial.LastWithdrawalAmount.String, Day: historial.WithdrawalDay, Time: historial.WithdrawalTime}
		movimientos = append(movimientos, mov)
	}
	views.GetMovements(movimientos).Render(h.ctx, w)
}
