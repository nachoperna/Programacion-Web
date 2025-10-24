package handlers

import (
	sqlc "TP_Especial/db/sqlc"
	"TP_Especial/views"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

type Request struct {
	FromAlias string `json:"from_alias"`
	ToAlias   string `json:"to_alias"`
	Amount    string `json:"amount"`
	Message   string `json:"message"`
}

func (h *Handler) ListRequestsTo(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("to_alias")
	pedidos, err := h.queries.GetRequestsTo(h.ctx, sqlc.GetRequestsToParams{
		ToAlias:   alias,
		SortBy:    "amount",
		SortOrder: "desc",
	})
	if len(pedidos) == 0 {
		views.SinPedidos(true).Render(h.ctx, w)
		return
	}
	if err != nil {
		http.Error(w, "Error obteniendo requests", http.StatusNotFound)
		return
	}
	views.GetRequestsTo(pedidos).Render(h.ctx, w)
}

func (h *Handler) ListRequestsFrom(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("from_alias")
	pedidos, err := h.queries.GetRequestsFrom(h.ctx, sqlc.GetRequestsFromParams{
		FromAlias: alias,
		SortBy:    "amount",
		SortOrder: "desc",
	})
	if len(pedidos) == 0 {
		views.SinPedidos(false).Render(h.ctx, w)
		return
	}
	if err != nil {
		http.Error(w, "Error obteniendo requests", http.StatusNotFound)
		return
	}
	views.GetRequestsFrom(pedidos).Render(h.ctx, w)
}

func (h *Handler) DeleteRequestsTo(w http.ResponseWriter, r *http.Request) {
	from_alias := r.URL.Query().Get("from_alias")
	to_alias := r.URL.Query().Get("to_alias")
	_, err := h.queries.DeleteRequest(h.ctx, sqlc.DeleteRequestParams{
		FromAlias: from_alias,
		ToAlias:   to_alias,
	})
	if err != nil {
		http.Error(w, "Error obteniendo requests", http.StatusNotFound)
		return
	}
	h.ListRequestsTo(w, r)
}

func (h *Handler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	// Decodifica el arreglo de objetos JSON del cuerpo de la petición
	var requests []Request
	err := json.NewDecoder(r.Body).Decode(&requests)
	if err != nil {
		http.Error(w, "Error: El JSON enviado es inválido o no es un arreglo.", http.StatusBadRequest)
		return
	}

	// Prepara slices para llevar un registro de las operaciones
	var creados []sqlc.MoneyRequest
	fallidos := make(map[string]string)

	// Itera sobre cada pedido recibido
	for _, req := range requests {
		// Crea el pedido en la base de datos
		creado, err := h.queries.InsertRequest(h.ctx, sqlc.InsertRequestParams{
			FromAlias: req.FromAlias,
			ToAlias:   req.ToAlias,
			Amount:    req.Amount,
			Message:   sql.NullString{String: req.Message, Valid: true}, // Pasamos el sql.NullString directamente
		})

		if err != nil {
			// Si hay un error, lo registramos
			errorKey := fmt.Sprintf("De '%s' a '%s'", req.FromAlias, req.ToAlias)
			fallidos[errorKey] = err.Error()
		} else {
			// Si es exitoso, lo añadimos a la lista de creados
			creados = append(creados, creado)
		}
	}

	// Prepara la respuesta JSON
	respuesta := map[string]any{
		"pedidos_creados":  creados,
		"pedidos_fallidos": fallidos,
	}

	// Determina el código de estado HTTP adecuado
	status := http.StatusOK
	if len(creados) > 0 && len(fallidos) == 0 {
		status = http.StatusCreated // Todo salió bien
	} else if len(creados) == 0 && len(fallidos) > 0 {
		status = http.StatusBadRequest // Ninguno se pudo crear
	} else if len(creados) > 0 && len(fallidos) > 0 {
		status = http.StatusMultiStatus // Algunos salieron bien, otros no
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(respuesta)
}
