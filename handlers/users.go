package handlers

import (
	sqlc "TP_Especial/db/sqlc"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"net/http"
)

type User struct {
	Alias    string `json:"alias"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	users, err := h.queries.ListUsers(h.ctx)
	if err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	var creados []sqlc.User
	fallidos := make(map[string]string)

	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Error: El JSON enviado es inválido.", http.StatusBadRequest)
		return
	}
	for _, user := range users {
		_, err := h.queries.GetUser(h.ctx, user.Alias)
		if err == sql.ErrNoRows {
			creado, err := h.queries.InsertUser(h.ctx, sqlc.InsertUserParams{
				Alias:    user.Alias,
				Name:     user.Name,
				Email:    user.Email,
				Password: user.Password})
			if err != nil {
				fallidos[user.Alias] = "Error al insertar usuario"
			} else {
				creados = append(creados, creado)
			}
		} else {
			fallidos[user.Alias] = "Usuario ya registrado"
		}
	}

	respuesta := map[string]any{
		"usuarios_creados":  creados,
		"usuarios_fallidos": fallidos,
	}

	status := http.StatusOK
	if len(creados) == 0 && len(fallidos) > 0 {
		status = http.StatusBadRequest
	} else if len(creados) > 0 && len(fallidos) == 0 {
		status = http.StatusCreated
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(respuesta)
}

func (h *Handler) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	err := h.queries.DeleteAllUsers(h.ctx)
	if err != nil {
		http.Error(w, "Error eliminando todos los usuarios", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("alias")
	user, err := h.queries.GetUser(h.ctx, alias)
	w.Header().Set("Content-type", "application/json")
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("El usuario no se encuentra registrado en el sistema")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("alias")

	_, err := h.queries.DeleteUser(h.ctx, alias)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("El usuario no se encuentra registrado en el sistema")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Error eliminando usuario.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var datos User
	alias := r.URL.Query().Get("alias")

	err := json.NewDecoder(r.Body).Decode(&datos)
	if err != nil {
		http.Error(w, "Error: El JSON enviado es inválido.", http.StatusBadRequest)
		return
	}
	_, err = h.queries.UpdateUser(h.ctx, sqlc.UpdateUserParams{
		Alias:    alias,
		Name:     datos.Name,
		Email:    datos.Email,
		Password: datos.Password,
	})
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("El usuario no se encuentra registrado en el sistema")
		return
	} else if err != nil {
		http.Error(w, "Error actualizando usuario.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
