package main

import (
	sqlc "TP_Especial/db/sqlc"
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

var newUsers []sqlc.InsertUserParams

func TestQueries_CRUD(test *testing.T) {
	datos_conexion := "host=localhost port=5432 user=nachoperna password=nachobdtpe dbname=BD_TPEspecial sslmode=disable"
	db, err := sql.Open("postgres", datos_conexion)
	if err != nil {
		test.Fatalf("Error al conectar a la Base: %v", err)
	}

	queries := sqlc.New(db)

	test.Run("Creating Users and Accounts", func(t *testing.T) {
		tCreateUser(queries, t)
	})

	test.Run("Getting Users", func(t *testing.T) {
		tGetUser(queries, t)
	})

	test.Run("Updating Users", func(t *testing.T) {
		tUpdateUser(queries, t)
	})

	test.Run("Listing Users and Accounts", func(t *testing.T) {
		tListUsers(queries, t)
	})

	test.Run("Deleting Users and Accounts", func(t *testing.T) {
		tDeleteUser(queries, t)
	})
}

func match(user1 sqlc.InsertUserParams, user2 sqlc.InsertUserParams) bool {
	return user1.Alias == user2.Alias
}

func tCreateUser(queries *sqlc.Queries, t *testing.T) {
	user1 := sqlc.InsertUserParams{Alias: "aliasu1", Name: "user1", Email: "user1@mail.com", Password: "u1pass"}
	user2 := sqlc.InsertUserParams{Alias: "aliasu2", Name: "user2", Email: "user2@mail.com", Password: "u2pass"}
	user3 := sqlc.InsertUserParams{Alias: "aliasu3", Name: "user3", Email: "user3@mail.com", Password: "u3pass"}

	err := queries.InsertUser(t.Context(), user1)
	if err != nil {
		t.Fatalf("Error creando Usuario 1: %v", err)
	}
	newUsers = append(newUsers, user1)
	acc_get, err := queries.GetAccount(t.Context(), user1.Alias)
	if err == sql.ErrNoRows {
		t.Fatalf("No se creó la cuenta con alias: %v", user1.Alias)
	}
	fmt.Printf("\n-Se creó cuenta con alias %s y balance: $%s", acc_get.Alias, acc_get.Balance)

	err = queries.InsertUser(t.Context(), user2)
	if err != nil {
		t.Fatalf("Error creando Usuario 2: %v", err)
	}
	newUsers = append(newUsers, user2)
	acc_get, err = queries.GetAccount(t.Context(), user2.Alias)
	if err == sql.ErrNoRows {
		t.Fatalf("No se creó la cuenta con alias: %v", user2.Alias)
	}
	fmt.Printf("\n-Se creó cuenta con alias %s y balance: $%s", acc_get.Alias, acc_get.Balance)

	err = queries.InsertUser(t.Context(), user3)
	if err != nil {
		t.Fatalf("Error creando Usuario 3: %v", err)
	}
	newUsers = append(newUsers, user3)
	acc_get, err = queries.GetAccount(t.Context(), user3.Alias)
	if err == sql.ErrNoRows {
		t.Fatalf("No se creó la cuenta con alias: %v", user3.Alias)
	}
	fmt.Printf("\n-Se creó cuenta con alias %s y balance: $%s", acc_get.Alias, acc_get.Balance)
}

func tGetUser(queries *sqlc.Queries, t *testing.T) {
	for _, u := range newUsers {
		u_get, err := queries.GetUser(t.Context(), u.Alias)
		if err != nil {
			t.Fatalf("Error obteniendo usuario con id: %v", u.Alias)
		}
		if match(u, sqlc.InsertUserParams{Alias: u_get.Alias}) {
			fmt.Println("El usuario obtenido y el creaado son iguales.")
		} else {
			t.Fatalf("Error obteniendo usuario, el buscado y el obtenido no son iguales: %v", u_get.Alias)
		}
	}
}

func tUpdateUser(queries *sqlc.Queries, t *testing.T) {
	for _, u := range newUsers {
		err := queries.UpdateUser(t.Context(), sqlc.UpdateUserParams{Alias: u.Alias, Name: u.Name, Email: u.Email + "[ACTUALIZADO]"})
		if err != nil {
			t.Fatalf("Error actualizando usuario con Alias: %v", u.Alias)
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
			fmt.Printf("- Usuario [Alias: %s, Name: %s, Email: %s, UltSession: %s]\n", u.Alias, u.Name, u.Email, u.SignedUp)
			account, err := queries.GetAccount(t.Context(), u.Alias)
			if err != nil {
				t.Fatalf("Error obteniendo cuenta: %v", err)
			}
			fmt.Printf("-- Cuenta [Alias: %s, Balance: $%s]", account.Alias, account.Balance)
		} else {
			fmt.Printf("El usuario con Alias: %s no se encuentra actualizado en la Base", u.Alias)
		}
	}
}

func tDeleteUser(queries *sqlc.Queries, t *testing.T) {
	for _, u := range newUsers {
		err := queries.DeleteUser(t.Context(), u.Alias)
		if err != nil {
			t.Fatalf("Error eliminando usuario con Alias: %s", u.Alias)
		}
		_, err = queries.GetUser(t.Context(), u.Alias)
		if err == nil {
			t.Fatalf("El usuario con Alias: %s se 'elimino' pero se sigue pudiendo obtener", u.Alias)
		} else if err == sql.ErrNoRows {
			fmt.Printf("- El usuario con alias: %s ya no existe en la Base\n", u.Alias)
			_, err := queries.GetAccount(t.Context(), u.Alias)
			if err == sql.ErrNoRows {
				fmt.Printf("-- La cuenta con alias: %s ya no existe en la Base\n", u.Alias)
			}
		}
	}
}

