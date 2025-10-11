package main

import (
	// "context"
	"database/sql"
	"fmt"
	"log"

	// sqlc "SQLC/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	datos_conexion := "user=nachoperna password=nachobdtp2sqlc dbname=BD_tp2sqlc"
	db, err := sql.Open("postgres", datos_conexion)
	if err != nil {
		log.Fatalf("Error al conectar con la Base de Datos: %v", err)
	}
	defer db.Close()

	// Verificar que la conexión funciona
	err = db.Ping()
	if err != nil {
		log.Fatal("Error al verificar conexión:", err)
	}

	fmt.Println("✅ Conectado a PostgreSQL exitosamente!")
}
