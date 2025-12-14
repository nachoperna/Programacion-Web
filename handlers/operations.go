package handlers

import (
	sqlc "TP_Especial/db/sqlc"
	"TP_Especial/views"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias":  r.FormValue("alias"),
		"Amount": r.FormValue("amount"),
	}
	amount, err := strconv.ParseFloat(datos["Amount"], 64)
	if err != nil {
		http.Error(w, "Monto inválido", http.StatusBadRequest)
		return
	}

	info, valido := h.DepositLogic(w, datos["Alias"], amount, "Deposito")
	if valido {
		datos["Email"] = r.URL.Query().Get("email")
		views.SetInfo(datos["Alias"], h.balanceFormateado(info.Balance), datos["Email"], info.LastMovementType.String).Render(h.ctx, w)
	}
}

func (h *Handler) DepositLogic(w http.ResponseWriter, alias string, amount float64, operation_type string) (sqlc.DepositRow, bool) {
	if amount <= 0 {
		w.Header().Set("HX-Trigger", "invalid_amount")
		w.WriteHeader(http.StatusBadRequest)
		return sqlc.DepositRow{}, false
	}
	info, err := h.queries.Deposit(h.ctx, sqlc.DepositParams{
		Alias:              alias,
		Balance:            fmt.Sprintf("%.2f", amount),
		LastMovementType:   sql.NullString{String: operation_type, Valid: true},
		LastMovementAmount: sql.NullString{String: fmt.Sprintf("%.2f", amount), Valid: true},
	})
	if err != nil {
		http.Error(w, "Error al Depositar", http.StatusInternalServerError)
		return sqlc.DepositRow{}, false
	}

	return info, true
}

func (h *Handler) Withdrawal(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias":  r.FormValue("alias"),
		"Amount": r.FormValue("amount"),
	}
	amount, err := strconv.ParseFloat(datos["Amount"], 64)
	if err != nil {
		http.Error(w, "Monto inválido", http.StatusBadRequest)
		return
	}

	info, valido := h.WithdrawalLogic(w, datos["Alias"], amount, "Retiro")
	if valido {
		datos["Email"] = r.URL.Query().Get("email")
		views.SetInfo(datos["Alias"], h.balanceFormateado(info.Balance), datos["Email"], info.LastMovementType.String).Render(h.ctx, w)
	}
}

func (h *Handler) WithdrawalLogic(w http.ResponseWriter, alias string, amount float64, operation_type string) (sqlc.WithdrawalRow, bool) {
	if amount <= 0 {
		w.Header().Set("HX-Trigger", "invalid_amount")
		w.WriteHeader(http.StatusBadRequest)
		return sqlc.WithdrawalRow{}, false
	}
	if !h.EnoughBalance(w, alias, amount) {
		w.Header().Set("HX-Trigger", "not_enough_balance")
		w.WriteHeader(http.StatusBadRequest)
		return sqlc.WithdrawalRow{}, false
	}
	info, err := h.queries.Withdrawal(h.ctx, sqlc.WithdrawalParams{
		Alias:              alias,
		Balance:            fmt.Sprintf("%.2f", amount),
		LastMovementType:   sql.NullString{String: operation_type, Valid: true},
		LastMovementAmount: sql.NullString{String: fmt.Sprintf("%.2f", amount), Valid: true},
	})
	if err != nil {
		http.Error(w, "Error al quitar dinero de la cuenta", http.StatusInternalServerError)
		return sqlc.WithdrawalRow{}, false
	}
	return info, true
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias_propio": r.FormValue("own_alias"),
		"Alias_otro":   r.FormValue("other_alias"),
		"Amount":       r.FormValue("amount"),
	}

	amount, err := strconv.ParseFloat(datos["Amount"], 64)
	if err != nil {
		w.Header().Set("HX-Trigger", "invalid_amount")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	info, valido := h.TransferLogic(w, datos["Alias_propio"], datos["Alias_otro"], amount)
	if valido {
		datos["Email"] = r.URL.Query().Get("email")
		views.SetInfo(datos["Alias_propio"], h.balanceFormateado(info.Balance), datos["Email"], info.LastMovementType.String).Render(h.ctx, w)
	}
}

func (h *Handler) TransferLogic(w http.ResponseWriter, own_alias, other_alias string, amount float64) (sqlc.WithdrawalRow, bool) {
	if own_alias == other_alias {
		w.Header().Set("HX-Trigger", "mismo_alias")
		w.WriteHeader(http.StatusBadRequest)
		return sqlc.WithdrawalRow{}, false
	}

	_, err := h.queries.GetUser(h.ctx, other_alias)
	if err == sql.ErrNoRows {
		w.Header().Set("HX-Trigger", "alias_not_found")
		w.WriteHeader(http.StatusBadRequest)
		return sqlc.WithdrawalRow{}, false
	}

	info, valido1 := h.WithdrawalLogic(w, own_alias, amount, "Transferencia")
	_, valido2 := h.DepositLogic(w, other_alias, amount, "Transferencia")

	if !valido1 || !valido2 {
		return sqlc.WithdrawalRow{}, false
	}
	return info, true
}

func (h *Handler) RequestMoney(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias_propio": r.FormValue("from_alias"),
		"Alias_otro":   r.FormValue("to_alias"),
		"Amount":       r.FormValue("amount"),
		"Mensaje":      r.FormValue("message"),
	}

	if datos["Alias_propio"] == datos["Alias_otro"] {
		w.Header().Set("HX-Trigger", "mismo_alias")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(datos["Amount"], 64)
	if err != nil || amount <= 0 {
		w.Header().Set("HX-Trigger", "invalid_amount")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.queries.InsertRequest(h.ctx, sqlc.InsertRequestParams{
		FromAlias: datos["Alias_propio"],
		ToAlias:   datos["Alias_otro"],
		Amount:    datos["Amount"],
		Message:   sql.NullString{String: datos["Mensaje"], Valid: true},
	})

	_ = h.queries.UpdateHistory(h.ctx, sqlc.UpdateHistoryParams{
		Alias:  datos["Alias_propio"],
		Type:   "Pedido",
		Amount: fmt.Sprintf("%.2f", amount),
	})
	redirectURL := fmt.Sprintf("/home?alias=%s",
		datos["Alias_propio"])

	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) EnoughBalance(w http.ResponseWriter, alias string, monto float64) bool {
	balance, err := h.queries.GetBalance(h.ctx, alias)
	if err != nil {
		http.Error(w, "Error al obtener balance", http.StatusInternalServerError)
	}
	balanceS := balance.Balance                     // Se obtiene el valor String balance de la fila devuelta
	balanceP, _ := strconv.ParseFloat(balanceS, 64) // Se parsea el string a float64
	if monto > balanceP {
		return false
	}
	return true
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := h.queries.GetBalance(h.ctx, r.URL.Query().Get("user"))
	if err != nil {
		http.Error(w, "Error al obtener balance", http.StatusInternalServerError)
		return
	}
	io.WriteString(w, balance.Balance)
}
