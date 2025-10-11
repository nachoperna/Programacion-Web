package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Servidor de archivos
	fs := http.FileServer(http.Dir("static/"))

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

	// Procesamiento del formulario
	http.HandleFunc("/contacto", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Parsear el formulario
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
			return
		}

		datos := map[string]string{
			"Nombre":  r.FormValue("user"),
			"Email":   r.FormValue("email"),
			"Mensaje": r.FormValue("message"),
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
