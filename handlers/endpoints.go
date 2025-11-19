package handlers

import (
	_ "github.com/lib/pq"
	"net/http"
)

func (h *Handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUsers(w, r)
	case http.MethodPost:
		h.CreateUsers(w, r)
	case http.MethodDelete:
		h.DeleteUsers(w, r)
	default:
		http.Error(w, "Metodo invalido", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUser(w, r)
	case http.MethodPost:
		h.CreateUsers(w, r)
	case http.MethodPut:
		h.UpdateUser(w, r)
	case http.MethodDelete:
		h.DeleteUser(w, r)
	default:
		http.Error(w, "Metodo invalido", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) RequestsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		h.CreateRequest(w, r)
	case http.MethodDelete:
	default:
		http.Error(w, "Metodo invalido", http.StatusMethodNotAllowed)
	}
}
