package handlers

import (
	sqlc "TP_Especial/db/sqlc"
	"TP_Especial/views"
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (h *Handler) ConfirmLogin(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias":        r.FormValue("alias"),
		"Password":     r.FormValue("password"),
		"New_Password": r.FormValue("new-pass"),
	}
	var user sqlc.User
	var err error

	if datos["New_Password"] != "" {
		user, err = h.queries.UpdateUser(h.ctx, sqlc.UpdateUserParams{
			Alias:    datos["Alias"],
			Password: datos["New_Password"],
		})
		if err == sql.ErrNoRows {
			w.Header().Set("HX-Trigger", "alias_not_found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		user, err = h.queries.GetUser(h.ctx, datos["Alias"])
		if err == sql.ErrNoRows {
			w.Header().Set("HX-Trigger", "alias_not_found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if user.Password != datos["Password"] {
			w.Header().Set("HX-Trigger", "password_incorrect")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	redirectURL := fmt.Sprintf("/home?alias=%s&name=%s&email=%s",
		user.Alias,
		user.Name,
		user.Email)

	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ConfirmRegiser(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias":    r.FormValue("alias"),
		"Name":     r.FormValue("name"),
		"Email":    r.FormValue("email"),
		"Password": r.FormValue("password"),
	}

	_, err := h.queries.GetUser(h.ctx, datos["Alias"])
	if err == sql.ErrNoRows {
		_, err = h.queries.InsertUser(h.ctx, sqlc.InsertUserParams{
			Alias:    datos["Alias"],
			Name:     datos["Name"],
			Email:    datos["Email"],
			Password: datos["Password"]})
		if err != nil {
			w.Header().Set("HX-Trigger", "error_registro")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		w.Header().Set("HX-Trigger", "alias_usado")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	redirectURL := fmt.Sprintf("/home?alias=%s&name=%s&email=%s",
		datos["Alias"],
		datos["Name"],
		datos["Email"])

	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Showhome(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	balance, err := h.queries.GetBalance(h.ctx, values.Get("alias"))
	if err != nil {
		http.Error(w, "Error al obtener datos de cuenta", http.StatusInternalServerError)
		return
	}

	// libreria para poder formatear un numero al formato español de puntos para miles y comas para decimales $100.000,00
	printer := message.NewPrinter(language.Spanish)
	balance_float, _ := strconv.ParseFloat(balance.Balance, 64)
	balance_formateado := printer.Sprintf("%.2f", balance_float)
	datos := map[string]any{
		"Alias":            values.Get("alias"),
		"Balance":          balance_formateado,
		"LastMovementType": balance.LastMovementType.String,
	}

	if values.Get("name") == "" { // significa que llegamos a home a traves de una operacion y debemos ir a la base a obtener los datos faltantes
		user, _ := h.queries.GetUser(h.ctx, values.Get("alias"))
		datos["Name"] = user.Name
		datos["Email"] = user.Email
	} else {
		datos["Name"] = values.Get("name")
		datos["Email"] = values.Get("email")
	}

	// Servir el template con datos actualizados
	pedidos, _ := h.queries.GetRequestsToCount(h.ctx, values.Get("alias"))
	buf := new(bytes.Buffer)
	views.RequestCount(int(pedidos), datos["Alias"].(string), datos["Email"].(string)).Render(h.ctx, buf)
	datos["BotonPedidos"] = template.HTML(buf.String())

	tmp, err := template.ParseFiles("static/bienvenida.html")
	if err != nil {
		http.Error(w, "Error al cargar template", http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, datos)
}
