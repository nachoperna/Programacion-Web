package main

import (
	sqlc "TP_Especial/db/sqlc"
	"context"
	"database/sql"
	"encoding/json"
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

type User struct {
	Alias    string `json:"alias"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

	// Endpoint REST para tabla Users
	http.HandleFunc("/users", usersHandler)

	// Endpoint REST para tabla Users
	http.HandleFunc("/users/", userHandler)

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
		"Alias":    r.FormValue("alias"),
		"Password": r.FormValue("password"),
	}

	user, err := queries.GetUser(ctx, datos["Alias"])
	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/?error=alias_not_found", http.StatusSeeOther)
		return
	}
	if user.Password != datos["Password"] {
		http.Redirect(w, r, "/?error=password_incorrect", http.StatusSeeOther)
		return
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
			http.Error(w, "Error al insertar usuario en la base", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "El alias ya se encuentra registrado en la Base. Por favor elija otro", http.StatusInternalServerError)
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
	// POR EL MOMMENTO EL USUARIO DEBE INGRESAR SU PROPIO ALIAS
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
	// POR EL MOMMENTO EL USUARIO DEBE INGRESAR SU PROPIO ALIAS
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

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		createUsers(w, r)
	case http.MethodDelete:
		deleteUsers(w, r)
	default:
		http.Error(w, "Metodo invalido", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUser(w, r)
	case http.MethodPost:
		createUsers(w, r)
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, "Metodo invalido", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	users, err := queries.ListUsers(ctx)
	if err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func createUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	var creados []sqlc.User
	fallidos := make(map[string]string)

	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Error: El JSON enviado es inválido.", http.StatusBadRequest)
		return
	}
	for _, user := range users {
		_, err := queries.GetUser(ctx, user.Alias)
		if err == sql.ErrNoRows {
			creado, err := queries.InsertUser(ctx, sqlc.InsertUserParams{
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

	respuesta := map[string]interface{}{
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

func deleteUsers(w http.ResponseWriter, r *http.Request) {
	err := queries.DeleteAllUsers(ctx)
	if err != nil {
		http.Error(w, "Error eliminando todos los usuarios", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("alias")
	user, err := queries.GetUser(ctx, alias)
	w.Header().Set("Content-type", "application/json")
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("El usuario no se encuentra registrado en el sistema")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("alias")

	_, err := queries.DeleteUser(ctx, alias)
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

func updateUser(w http.ResponseWriter, r *http.Request) {
	var datos User
	alias := r.URL.Query().Get("alias")

	err := json.NewDecoder(r.Body).Decode(&datos)
	if err != nil {
		http.Error(w, "Error: El JSON enviado es inválido.", http.StatusBadRequest)
		return
	}
	_, err = queries.UpdateUser(ctx, sqlc.UpdateUserParams{
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
