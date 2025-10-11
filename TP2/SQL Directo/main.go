package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

type User struct {
	id    int
	name  string
	email string
}

// Servidor de archivos
var fs http.Handler = http.FileServer(http.Dir("static/"))

func main() {
	// fs := http.FileServer(http.Dir("static/"))
	conexion := "host=localhost port=5432 user=nachoperna password=nachobdtp2 dbname=BD_tp2 sslmode=disable"

	// Conectar a la base de datos
	db, err := sql.Open("postgres", conexion)
	if err != nil {
		log.Fatal("Error al conectar:", err)
	}
	defer db.Close() // Cerrar conexión al final

	// Verificar que la conexión funciona
	err = db.Ping()
	if err != nil {
		log.Fatal("Error al verificar conexión:", err)
	}

	fmt.Println("✅ Conectado a PostgreSQL exitosamente!")

	// Pagina Principal
	http.HandleFunc("/", inicio)

	// Procesamiento del formulario
	http.HandleFunc("/contacto", procesarFormulario)

	http.HandleFunc("/borrarTablas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Solo método POST permitido", http.StatusMethodNotAllowed)
			return
		}

		// Función simple directa
		err := borrarTablas(db)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "✅ Todas las tablas fueron eliminadas")
		fmt.Fprint(w, "Se debe reiniciar el servidor")
	})

	fmt.Print("Servidor escuchando en puerto :8080")
	http.ListenAndServe(":8080", nil)
}

func inicio(w http.ResponseWriter, r *http.Request) {
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
}

func procesarFormulario(w http.ResponseWriter, r *http.Request) {
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
}

func borrarTablas(db *sql.DB) error {
	tables, err := getAllTables(db)
	// Borrar cada tabla
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			log.Fatal(err)
		}
	}
	return err
}

func getAllTables(db *sql.DB) ([]string, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
		ORDER BY table_name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, rows.Err()
}

func (r *userRepository) CreateUser(user *User) error {
	_, err := r.db.Exec("insert into users (id, name, email) values ($1, $2, $3)", user.id, user.name, user.email)
	return err
}

func (r *userRepository) GetUserByID(id int) (*User, error) {
	row := r.db.QueryRow("select id, name, email from users where id = $1", id)
	user := &User{}
	err := row.Scan(&user.id, &user.name, &user.email)
	if err == sql.ErrNoRows {
		fmt.Println("No se encontro User con el id: ", id)
		return nil, err
	}
	return user, err
}

func (r *userRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("delete from users where id = $1", id)
	return err
}

func (r *userRepository) ListUsers() ([]*User, error) {
	rows, err := r.db.Query("select id, name, email from users")
	if err != nil {
		return nil, fmt.Errorf("Error al buscar usuarios: %v", err)
	}
	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.id, &user.name, &user.email)
		if err != nil {
			return nil, fmt.Errorf("Error iterando filas: %v", err)
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepository) UpdateUser(user *User) error {
	_, err := r.db.Exec("update users set name = $1, email = $2 where id = $3", user.name, user.email, user.id)
	return err
}

// func showUsers(users []*User) {
// 	if len(users) == 0 {
// 		fmt.Println("No hay usuarios para mostrar")
// 		return
// 	}
// 	fmt.Printf("Hay un total de %d usuarios a mostrar\n", len(users))
// 	for _, user := range users {
// 		fmt.Printf("ID: %d, name: %s, email: %s", user.id, user.name, user.email)
// 	}
// }
