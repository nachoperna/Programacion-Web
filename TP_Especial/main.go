package main

import (
	sqlc "TP_Especial/db/sqlc"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/lib/pq"
)

var queries *sqlc.Queries
var ctx context.Context

func main() {
	// Servidor de archivos
	fs := http.FileServer(http.Dir("static/"))

	// Conexion a la Base de Datos PostgreSQL
	db, err := sql.Open("postgres", "host=localhost port=5432 user=nachoperna password=nachobdtpe dbname=BD_TPEspecial sslmode=disable")
	if err != nil {
		log.Fatalf("Error al conectar con la Base de Datos: %v", err)
	}
	defer db.Close()
	fmt.Println("✅ Conectado a PostgreSQL correctmente")

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar si el archivo existe físicamente
		requestedPath := filepath.Join("static", r.URL.Path)
		_, err := os.Stat(requestedPath)
		if os.IsNotExist(err) {
			// Archivo no existe, servir página personalizada
			http.ServeFile(w, r, "static/ruta_invalida.html")
			return
		}
		// Si existe, usar FileServer normalmente
		fs.ServeHTTP(w, r)
	}))

	queries = sqlc.New(db)
	ctx = context.Background()

	// Procesamiento del formulario de LOGIN
	http.HandleFunc("/login-confirmation", confirmLogin)

	// Procesamiento del formulario de REGISTRO
	http.HandleFunc("/register-confirmation", confirmRegiser)

	// Servir home del usuario
	http.HandleFunc("/home", showhome)

	// Procesamiento de DEPOSITO
	http.HandleFunc("/deposit", deposit)

	// Procesamiento de RETIRO
	http.HandleFunc("/withdrawal", withdrawal)

	// Procesamiento de TRANSFERENCIA
	http.HandleFunc("/transfer", transfer)

	fmt.Print("Servidor escuchando en puerto :8080")
	http.ListenAndServe(":8080", nil)
}

func showhome(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	balance, err := queries.GetBalance(ctx, values.Get("alias"))
	if err != nil {
		http.Error(w, "Error al obtener datos de cuenta", http.StatusInternalServerError)
		return
	}
	datos := map[string]string{
		"Alias":            values.Get("alias"),
		"Balance":          balance.Balance,
		"LastMovementType": balance.LastMovementType.String,
	}

	if values.Get("name") == "" { // significa que llegamos a home a traves de una operacion y debemos ir a la base a obtener los datos faltantes
		user, _ := queries.GetUser(ctx, values.Get("alias"))
		datos["Name"] = user.Name
		datos["Email"] = user.Email
	} else {
		datos["Name"] = values.Get("name")
		datos["Email"] = values.Get("email")
	}

	// Servir el template con datos actualizados
	tmp, err := template.ParseFiles("static/bienvenida.html")
	if err != nil {
		http.Error(w, "Error al cargar template", http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, datos)
}

func confirmLogin(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias":        r.FormValue("alias"),
		"Password":     r.FormValue("password"),
		"New_Password": r.FormValue("new-pass"),
	}
	var user sqlc.User
	var err error
	if datos["New_Password"] != "" {
		user, err = queries.UpdateUser(ctx, sqlc.UpdateUserParams{
			Alias:    datos["Alias"],
			Password: datos["New_Password"],
		})
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/?error=alias_not_found", http.StatusSeeOther)
			return
		}
	} else {
		user, err = queries.GetUser(ctx, datos["Alias"])
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/?error=alias_not_found", http.StatusSeeOther)
			return
		}
		if user.Password != datos["Password"] {
			http.Redirect(w, r, "/?error=password_incorrect", http.StatusSeeOther)
			return
		}
	}

	redirectURL := fmt.Sprintf("/home?alias=%s&name=%s&email=%s",
		user.Alias,
		user.Name,
		user.Email)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func confirmRegiser(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias":    r.FormValue("alias"),
		"Name":     r.FormValue("name"),
		"Email":    r.FormValue("email"),
		"Password": r.FormValue("password"),
	}

	_, err := queries.GetUser(ctx, datos["Alias"])
	if err == sql.ErrNoRows {
		_, err = queries.InsertUser(ctx, sqlc.InsertUserParams{
			Alias:    datos["Alias"],
			Name:     datos["Name"],
			Email:    datos["Email"],
			Password: datos["Password"]})
		if err != nil {
			http.Redirect(w, r, "/?error=error_registro", http.StatusSeeOther)
			return
		}
	} else {
		http.Redirect(w, r, "/?error=alias_usado", http.StatusSeeOther)
		return
	}

	redirectURL := fmt.Sprintf("/home?alias=%s&name=%s&email=%s",
		datos["Alias"],
		datos["Name"],
		datos["Email"])

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func deposit(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias":  r.FormValue("alias"),
		"Amount": r.FormValue("amount"),
	}
	amount, err := strconv.ParseFloat(datos["Amount"], 64)
	if err != nil {
		http.Error(w, "Monto inválido", http.StatusBadRequest)
		return
	}

	_, err = queries.Deposit(ctx, sqlc.DepositParams{
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

func transfer(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias_propio": r.FormValue("own_alias"),
		"Alias_otro":   r.FormValue("other_alias"),
		"Amount":       r.FormValue("amount"),
	}
	amount, err := strconv.ParseFloat(datos["Amount"], 64)
	if err != nil {
		http.Error(w, "Monto inválido", http.StatusBadRequest)
		return
	}

	_, err = queries.GetUser(ctx, datos["Alias_otro"])
	if err == sql.ErrNoRows {
		http.Error(w, "La alias destino no esta registrado en el sistema", http.StatusInternalServerError)
		return
	}
	err = queries.Transfer(ctx, sqlc.TransferParams{
		Alias:               datos["Alias_propio"],
		LastTransferAccount: sql.NullString{String: datos["Alias_otro"], Valid: true},
		LastTransferAmount:  sql.NullString{String: fmt.Sprintf("%.2f", amount), Valid: true},
	})
	if err != nil {
		fmt.Printf("error: %v", err)
		http.Error(w, "Error al quitar dinero de la cuenta de origen", http.StatusInternalServerError)
		return
	}

	_, err = queries.Deposit(ctx, sqlc.DepositParams{
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

func withdrawal(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"Alias":  r.FormValue("alias"),
		"Amount": r.FormValue("amount"),
	}
	amount, err := strconv.ParseFloat(datos["Amount"], 64)
	if err != nil {
		http.Error(w, "Monto inválido", http.StatusBadRequest)
		return
	}

	err = queries.Withdrawal(ctx, sqlc.WithdrawalParams{
		Alias:                datos["Alias"],
		LastWithdrawalAmount: sql.NullString{String: fmt.Sprintf("%.2f", amount), Valid: true},
	})
	if err != nil {
		http.Error(w, "Error al quitar dinero de la cuenta origen", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/home?alias=%s&name=%s&email=%s",
		datos["Alias"])

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// func newPassword(w http.ResponseWriter, r *http.Request){
//
// }
