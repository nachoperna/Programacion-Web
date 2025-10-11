package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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

	// Procesamiento del formulario de LOGIN
	http.HandleFunc("/login-confirmation", func(w http.ResponseWriter, r *http.Request) {

		datos := map[string]string{
			"Alias":    r.FormValue("alias"),
			"Password": r.FormValue("password"),
		}

		tmp, err := template.ParseFiles("static/bienvenida.html")
		if err != nil {
			http.Error(w, "Error al cargar datos", http.StatusInternalServerError)
			return
		}
		tmp.Execute(w, datos)
	})

	// Procesamiento del formulario de REGISTRO
	http.HandleFunc("/register-confirmation", func(w http.ResponseWriter, r *http.Request) {

		datos := map[string]string{
			"Alias":    r.FormValue("alias"),
			"Password": r.FormValue("password"),
		}

		tmp, err := template.ParseFiles("static/bienvenida.html")
		if err != nil {
			http.Error(w, "Error al cargar datos", http.StatusInternalServerError)
			return
		}
		tmp.Execute(w, datos)
	})
	fmt.Print("Servidor escuchando en puerto :8080")
	http.ListenAndServe(":8080", nil)
}
