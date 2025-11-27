package main

import (
	sqlc "TP_Especial/db/sqlc"
	"database/sql"
	"fmt"
	// "strings"
	"testing"
)

var newUsers []sqlc.InsertUserParams

func TestQueries_CRUD(test *testing.T) {
	datos_conexion := "host=localhost port=5432 user=nachoperna password=nachobdtpe dbname=BD_TPEspecial sslmode=disable"
	db, err := sql.Open("postgres", datos_conexion)
	if err != nil {
		test.Fatalf("Error al conectar a la Base: %v", err)
	}
	defer db.Close()
	h.queries.:= sqlc.New(db)

	test.Run("Creating Users and Accounts", func(t *testing.T) {
		tCreateUser(h.queries. t)
	})

	test.Run("Getting Users", func(t *testing.T) {
		tGetUser(h.queries. t)
	})

	test.Run("Updating Users", func(t *testing.T) {
		tUpdateUser(h.queries. t)
	})

	test.Run("Deposit", func(t *testing.T) {
		tDeposit(h.queries. t)
	})

	test.Run("Whitdrawal", func(t *testing.T) {
		tWhitdrawal(h.queries. t)
	})

	test.Run("Transfer", func(t *testing.T) {
		tTransfer(h.queries. t)
	})
	// test.Run("Listing Users and Accounts", func(t *testing.T) {
	// 	tListUsers(h.queries. t)
	// })
	//
	// test.Run("Deleting Users and Accounts", func(t *testing.T) {
	// 	tDeleteUser(h.queries. t)
	// })

	test.Run("Deleting All Users and Accounts", func(t *testing.T) {
		tDeleteAllUsers(h.queries. t)
	})

}

func match(user1 sqlc.InsertUserParams, user2 sqlc.InsertUserParams) bool {
	return user1.Alias == user2.Alias
}

func tCreateUser(h.queries.*sqlc.Queries, t *testing.T) {
	user1 := sqlc.InsertUserParams{Alias: "aliasu1", Name: "user1", Email: "user1@mail.com", Password: "u1pass"}
	user2 := sqlc.InsertUserParams{Alias: "aliasu2", Name: "user2", Email: "user2@mail.com", Password: "u2pass"}
	user3 := sqlc.InsertUserParams{Alias: "aliasu3", Name: "user3", Email: "user3@mail.com", Password: "u3pass"}

	err := h.queries.InsertUser(t.Context(), user1)
	if err != nil {
		t.Fatalf("Error creando Usuario 1: %v", err)
	}
	newUsers = append(newUsers, user1)
	acc_get, err := h.queries.GetAccount(t.Context(), user1.Alias)
	if err == sql.ErrNoRows {
		t.Fatalf("No se creó la cuenta con alias: %v", user1.Alias)
	}
	fmt.Printf("\n-Se creó cuenta con alias %s y balance: $%s", acc_get.Alias, acc_get.Balance)

	err = h.queries.InsertUser(t.Context(), user2)
	if err != nil {
		t.Fatalf("Error creando Usuario 2: %v", err)
	}
	newUsers = append(newUsers, user2)
	acc_get, err = h.queries.GetAccount(t.Context(), user2.Alias)
	if err == sql.ErrNoRows {
		t.Fatalf("No se creó la cuenta con alias: %v", user2.Alias)
	}
	fmt.Printf("\n-Se creó cuenta con alias %s y balance: $%s", acc_get.Alias, acc_get.Balance)

	err = h.queries.InsertUser(t.Context(), user3)
	if err != nil {
		t.Fatalf("Error creando Usuario 3: %v", err)
	}
	newUsers = append(newUsers, user3)
	acc_get, err = h.queries.GetAccount(t.Context(), user3.Alias)
	if err == sql.ErrNoRows {
		t.Fatalf("No se creó la cuenta con alias: %v", user3.Alias)
	}
	fmt.Printf("\n-Se creó cuenta con alias %s y balance: $%s", acc_get.Alias, acc_get.Balance)
}

func tGetUser(h.queries.*sqlc.Queries, t *testing.T) {
	for _, u := range newUsers {
		u_get, err := h.queries.GetUser(t.Context(), u.Alias)
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

func tUpdateUser(h.queries.*sqlc.Queries, t *testing.T) {
	for _, u := range newUsers {
		err := h.queries.UpdateUser(t.Context(), sqlc.UpdateUserParams{Alias: u.Alias, Name: u.Name, Email: u.Email + "[ACTUALIZADO]"})
		if err != nil {
			t.Fatalf("Error actualizando usuario con Alias: %v", u.Alias)
		}
	}
}

// func tListUsers(h.queries.*sqlc.Queries, t *testing.T) {
// 	users, err := h.queries.ListUsers(t.Context())
// 	if err != nil {
// 		t.Fatalf("Error listando usuarios: %v", err)
// 	}
// 	for _, u := range users {
// 		if strings.Contains(u.Email, "ACTUALIZADO") {
// 			fmt.Printf("- Usuario [Alias: %s, Name: %s, Email: %s, UltSession: %s]\n", u.Alias, u.Name, u.Email, u.SignedUp)
// 			account, err := h.queries.GetAccount(t.Context(), u.Alias)
// 			if err != nil {
// 				t.Fatalf("Error obteniendo cuenta: %v", err)
// 			}
// 			fmt.Printf("-- Cuenta [Alias: %s, Balance: $%s]", account.Alias, account.Balance)
// 		} else {
// 			fmt.Printf("El usuario con Alias: %s no se encuentra actualizado en la Base", u.Alias)
// 		}
// 	}
// }
//
// func tDeleteUser(h.queries.*sqlc.Queries, t *testing.T) {
// 	for _, u := range newUsers {
// 		err := h.queries.DeleteUser(t.Context(), u.Alias)
// 		if err != nil {
// 			t.Fatalf("Error eliminando usuario con Alias: %s", u.Alias)
// 		}
// 		_, err = h.queries.GetUser(t.Context(), u.Alias)
// 		if err == nil {
// 			t.Fatalf("El usuario con Alias: %s se 'elimino' pero se sigue pudiendo obtener", u.Alias)
// 		} else if err == sql.ErrNoRows {
// 			fmt.Printf("- El usuario con alias: %s ya no existe en la Base\n", u.Alias)
// 			_, err := h.queries.GetAccount(t.Context(), u.Alias)
// 			if err == sql.ErrNoRows {
// 				fmt.Printf("-- La cuenta con alias: %s ya no existe en la Base\n", u.Alias)
// 			}
// 		}
// 	}
// }

func tDeleteAllUsers(h.queries.*sqlc.Queries, t *testing.T) {
	err := h.queries.DeleteAllUsers(t.Context())
	if err != nil {
		t.Fatalf("Error eliminando todos los usuarios ")
	}
}

func tDeposit(h.queries.*sqlc.Queries, t *testing.T) {
	for _, user := range newUsers {
		err := h.queries.Deposit(t.Context(), sqlc.DepositParams{
			Alias:             user.Alias,
			LastDepositAmount: sql.NullString{String: fmt.Sprintf("%.2f", 100.00), Valid: true},
		})
		if err != nil {
			t.Fatalf("Error al depositar en usuario %s", user.Alias)
		}
		u, _ := h.queries.GetAccount(t.Context(), user.Alias)
		fmt.Printf("\n- Alias: %s, Balance: $%s", u.Alias, u.Balance)
	}
}

func tWhitdrawal(h.queries.*sqlc.Queries, t *testing.T) {
	for _, user := range newUsers {
		err := h.queries.Whitdrawal(t.Context(), sqlc.WhitdrawalParams{
			Alias:                user.Alias,
			LastWhitdrawalAmount: sql.NullString{String: fmt.Sprintf("%.2f", 50.00), Valid: true},
		})
		if err != nil {
			t.Fatalf("Error al retirar en usuario %s", user.Alias)
		}
		u, _ := h.queries.GetAccount(t.Context(), user.Alias)
		fmt.Printf("\n- Alias: %s, Balance: $%s", u.Alias, u.Balance)
	}
}

func tTransfer(h.queries.*sqlc.Queries, t *testing.T) {
	for _, user := range newUsers {
		err := h.queries.Transfer(t.Context(), sqlc.TransferParams{
			Alias:               user.Alias,
			LastTransferAccount: sql.NullString{String: newUsers[0].Alias, Valid: true},
			LastTransferAmount:  sql.NullString{String: fmt.Sprintf("%.2f", 50.00), Valid: true},
		})
		if err != nil {
			t.Fatalf("Error al retirar en usuario %s, cuenta origen", user.Alias)
		}
		err = h.queries.Deposit(t.Context(), sqlc.DepositParams{
			Alias:             newUsers[0].Alias,
			LastDepositAmount: sql.NullString{String: fmt.Sprintf("%.2f", 50.00), Valid: true},
		})
		if err != nil {
			t.Fatalf("Error al depositar dinero en la cuenta destino")
			return
		}
		origen, _ := h.queries.GetAccount(t.Context(), user.Alias)
		fmt.Printf("\n- Alias origen: %s, Balance origen: $%s", origen.Alias, origen.Balance)
		destino, _ := h.queries.GetAccount(t.Context(), newUsers[0].Alias)
		fmt.Printf("\n- Alias destino: %s, Balance destino: $%s", destino.Alias, destino.Balance)
	}
}
