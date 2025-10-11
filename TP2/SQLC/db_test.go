package main

import (
	sqlc "SQLC/db/sqlc"
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

var newUsers []sqlc.User

func TestQueries_CRUD(test *testing.T) {
	datos_conexion := "host=localhost port=5432 user=nachoperna password=nachobdtp2sqlc dbname=BD_tp2sqlc sslmode=disable"
	db, err := sql.Open("postgres", datos_conexion)
	if err != nil {
		test.Fatalf("Error al conectar a la Base: %v", err)
	}

	queries := sqlc.New(db)

	test.Run("Creating Users", func(t *testing.T) {
		tCreateUser(queries, t)
	})

	test.Run("Getting Users", func(t *testing.T) {
		tGetUser(queries, t)
	})

	test.Run("Updating Users", func(t *testing.T) {
		tUpdateUser(queries, t)
	})

	test.Run("Listing Users", func(t *testing.T) {
		tListUsers(queries, t)
	})

	test.Run("Deleting Users", func(t *testing.T) {
		tDeleteUser(queries, t)
	})

}

func match(user1 sqlc.User, user2 sqlc.CreateUserParams) bool {
	return user1.Name == user2.Name && user1.Email == user2.Email
}

func tCreateUser(queries *sqlc.Queries, t *testing.T) {
	user1 := sqlc.CreateUserParams{Name: "Sujero 1", Email: "sujero1@mail.com"}
	user2 := sqlc.CreateUserParams{Name: "Sujero 2", Email: "sujero2@mail.com"}
	user3 := sqlc.CreateUserParams{Name: "Sujero 3", Email: "sujero3@mail.com"}

	u, err := queries.CreateUser(t.Context(), user1)
	if err != nil {
		t.Fatalf("Error creando Usuario 1: %v", err)
	}
	newUsers = append(newUsers, u)
	if match(u, user1) {
		fmt.Println("Los datos del Usuario 1 coinciden")
	}

	u, err = queries.CreateUser(t.Context(), user2)
	if err != nil {
		t.Fatalf("Error creando Usuario 2: %v", err)
	}
	newUsers = append(newUsers, u)
	if match(u, user2) {
		fmt.Println("Los datos del Usuario 2 coinciden")
	}

	u, err = queries.CreateUser(t.Context(), user3)
	if err != nil {
		t.Fatalf("Error creando Usuario 3: %v", err)
	}
	newUsers = append(newUsers, u)
	if match(u, user3) {
		fmt.Println("Los datos del Usuario 3 coinciden")
	}
}

func tGetUser(queries *sqlc.Queries, t *testing.T) {
	for _, u := range newUsers {
		u_get, err := queries.GetUser(t.Context(), u.ID)
		if err != nil {
			t.Fatalf("Error obteniendo usuario con id: %v", u.ID)
		}
		if match(u, sqlc.CreateUserParams{Name: u_get.Name, Email: u_get.Email}) {
			fmt.Println("El usuario obtenido y el creaado son iguales.")
		}
	}
}

func tUpdateUser(queries *sqlc.Queries, t *testing.T) {
	for _, u := range newUsers {
		err := queries.UpdateUser(t.Context(), sqlc.UpdateUserParams{ID: u.ID, Name: u.Name, Email: u.Email + "[ACTUALIZADO"})
		if err != nil {
			t.Fatalf("Error actualizando usuario con id: ")
		}
	}
}

func tListUsers(queries *sqlc.Queries, t *testing.T) {
	users, err := queries.ListUsers(t.Context())
	if err != nil {
		t.Fatalf("Error listando usuarios: %v", err)
	}
	for _, u := range users {
		if strings.Contains(u.Email, "ACTUALIZADO") {
			fmt.Printf("Usuario [ID: %d, Name: %s, Email: %s]", u.ID, u.Name, u.Email)
		} else {
			fmt.Printf("El usuario con ID: %d no se encuentra actualizado en la Base", u.ID)
		}
	}
}

func tDeleteUser(queries *sqlc.Queries, t *testing.T) {
	for _, u := range newUsers {
		err := queries.DeleteUser(t.Context(), u.ID)
		if err != nil {
			t.Fatalf("Error eliminando usuario con ID: %d", u.ID)
		}
		_, err = queries.GetUser(t.Context(), u.ID)
		if err == nil {
			t.Fatalf("El usuario con ID: %d se 'elimino' pero se sigue pudiendo obtener", u.ID)
		}
		fmt.Printf("El usuario con ID: %d ya no existe en la Base", u.ID)
	}
}
