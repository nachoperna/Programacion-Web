package handlers

import (
	sqlc "TP_Especial/db/sqlc"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
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

	_, err = h.queries.Deposit(h.ctx, sqlc.DepositParams{
		Alias:             datos["Alias"],
		LastDepositAmount: sql.NullString{String: fmt.Sprintf("%.2f", amount), Valid: true},
	})

	if err != nil {
		http.Error(w, "Error al Depositar", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/home?alias=%s",
		datos["Alias"])

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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
	if !h.EnoughBalance(w, datos["Alias"], amount) {
		redirectURL := fmt.Sprintf("/home?alias=%s&error=not_enough_balance", datos["Alias"])
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	err = h.queries.Withdrawal(h.ctx, sqlc.WithdrawalParams{
		Alias:                datos["Alias"],
		LastWithdrawalAmount: sql.NullString{String: fmt.Sprintf("%.2f", amount), Valid: true},
	})
	if err != nil {
		http.Error(w, "Error al quitar dinero de la cuenta origen", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/home?alias=%s",
		datos["Alias"])

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias_propio": r.FormValue("own_alias"),
		"Alias_otro":   r.FormValue("other_alias"),
		"Amount":       r.FormValue("amount"),
	}
	if datos["Alias_propio"] == datos["Alias_otro"] {
		redirectURL := fmt.Sprintf("/home?alias=%s&error=mismo_alias", datos["Alias_propio"])
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(datos["Amount"], 64)
	if err != nil {
		http.Error(w, "Monto inválido", http.StatusBadRequest)
		return
	}

	if !h.EnoughBalance(w, datos["Alias_propio"], amount) {
		redirectURL := fmt.Sprintf("/home?alias=%s&error=not_enough_balance", datos["Alias_propio"])
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	_, err = h.queries.GetUser(h.ctx, datos["Alias_otro"])
	if err == sql.ErrNoRows {
		redirectURL := fmt.Sprintf("/home?alias=%s&error=alias_not_found", datos["Alias_propio"])
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	err = h.queries.Transfer(h.ctx, sqlc.TransferParams{
		Alias:               datos["Alias_propio"],
		LastTransferAccount: sql.NullString{String: datos["Alias_otro"], Valid: true},
		LastTransferAmount:  sql.NullString{String: fmt.Sprintf("%.2f", amount), Valid: true},
	})
	if err != nil {
		fmt.Printf("error: %v", err)
		http.Error(w, "Error al quitar dinero de la cuenta de origen", http.StatusInternalServerError)
		return
	}

	_, err = h.queries.Deposit(h.ctx, sqlc.DepositParams{
		Alias:             datos["Alias_otro"],
		LastDepositAmount: sql.NullString{String: fmt.Sprintf("%.2f", amount), Valid: true},
	})
	if err != nil {
		http.Error(w, "Error al depositar dinero en la cuenta destino", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/home?alias=%s",
		datos["Alias_propio"])

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *Handler) RequestMoney(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias_propio": r.FormValue("from_alias"),
		"Alias_otro":   r.FormValue("to_alias"),
		"Amount":       r.FormValue("amount"),
		"Mensaje":      r.FormValue("message"),
	}

	if datos["Alias_propio"] == datos["Alias_otro"] {
		redirectURL := fmt.Sprintf("/home?alias=%s&error=mismo_alias", datos["Alias_propio"])
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	amount, err := strconv.ParseFloat(datos["Amount"], 64)
	if err != nil || amount < 0 {
		redirectURL := fmt.Sprintf("/home?alias=%s&error=invalid_amount", datos["Alias_propio"])
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	_, err = h.queries.InsertRequest(h.ctx, sqlc.InsertRequestParams{
		FromAlias: datos["Alias_propio"],
		ToAlias:   datos["Alias_otro"],
		Amount:    datos["Amount"],
		Message:   sql.NullString{String: datos["Mensaje"], Valid: true},
	})

	redirectURL := fmt.Sprintf("/home?alias=%s",
		datos["Alias_propio"])

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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
