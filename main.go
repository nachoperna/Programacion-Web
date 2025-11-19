package main

import (
	sqlc "TP_Especial/db/sqlc"
	"TP_Especial/handlers"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Handler struct {
	queries *sqlc.Queries
	ctx     context.Context
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

	queries := sqlc.New(db)
	ctx := context.Background()

	handler := handlers.NewHandler(queries, ctx)

	// Procesamiento del formulario de LOGIN
	http.HandleFunc("/login-confirmation", handler.ConfirmLogin)

	// Procesamiento del formulario de REGISTRO
	http.HandleFunc("/register-confirmation", handler.ConfirmRegiser)

	// Servir home del usuario
	http.HandleFunc("/home", handler.Showhome)

	// Procesamiento de DEPOSITO
	http.HandleFunc("/deposit", handler.Deposit)

	// Procesamiento de RETIRO
	http.HandleFunc("/withdrawal", handler.Withdrawal)

	// Procesamiento de TRANSFERENCIA
	http.HandleFunc("/transfer", handler.Transfer)

	// Procesamiento de PEDIDOS DE DINERO
	http.HandleFunc("/moneyRequest", handler.RequestMoney)

	// Endpoint REST para tabla Users
	http.HandleFunc("/users", handler.UsersHandler)

	// Endpoint REST para tabla Users
	http.HandleFunc("/users/", handler.UserHandler)

	// Endpoint REST para tabla Money Requests
	http.HandleFunc("/requests", handler.RequestsHandler)

	// Procesamiento de TABLA DE PEDIDOS DE DINERO PARA MOSTRRAR EN HOME
	http.HandleFunc("/listRequestsTo", handler.ListRequestsTo)
	http.HandleFunc("/listRequestsFrom", handler.ListRequestsFrom)

	// Procesamiento de TABLA DE PEDIDOS DE DINERO PARA MOSTRRAR EN HOME
	http.HandleFunc("/deleteRequestsTo", handler.DeleteRequestsTo)

	fmt.Print("Servidor escuchando en puerto :8080")
	http.ListenAndServe(":8080", nil)
}
