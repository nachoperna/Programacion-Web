package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

func TestUserRepository_CRUD(test *testing.T) {
	conexion := "host=localhost port=5432 user=nachoperna password=nachobdtp2 dbname=BD_tp2 sslmode=disable"

	// Conectar a la base de datos
	db, err := sql.Open("postgres", conexion)
	if err != nil {
		log.Fatal("Error al conectar:", err)
	}
	defer db.Close() // Cerrar conexión al final
	fmt.Println("Testing conectado a la base correctamente")

	r := &userRepository{db: db}

	test.Run("Carga de usuarios", func(t *testing.T) {
		newUsers(r, test)
	})

	test.Run("Obtencion de usuarios", func(t *testing.T) {
		getUser(r, test)
	})

	test.Run("Actualizacion de usuarios", func(t *testing.T) {
		updateUser(r, test)
	})

	test.Run("Listado de usuarios", func(t *testing.T) {
		listUsers(r, test)
	})

	test.Run("Eliminacion de usuarios", func(t *testing.T) {
		deleteUser(r, test)
	})

}

func newUsers(r *userRepository, t *testing.T) {
	err := r.CreateUser(&User{1, "sujeto1", "sujero1@mail.com"})
	if err != nil {
		t.Errorf("Error creando usuario 1: %v", err)
	}

	err = r.CreateUser(&User{2, "sujeto2", "sujero2@mail.com"})
	if err != nil {
		t.Errorf("Error creando usuario 2: %v", err)
	}

	err = r.CreateUser(&User{3, "sujeto3", "sujero3@mail.com"})
	if err != nil {
		t.Errorf("Error creando usuario 3: %v", err)
	}
}

func getUser(r *userRepository, t *testing.T) {
	u, err := r.GetUserByID(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Usuario encontrado con datos: [id: %d] [name: %s] [email: %s]", u.id, u.name, u.email)
	u, err = r.GetUserByID(2)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Usuario encontrado con datos: [id: %d] [name: %s] [email: %s]", u.id, u.name, u.email)

	u, err = r.GetUserByID(3)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Usuario encontrado con datos: [id: %d] [name: %s] [email: %s]", u.id, u.name, u.email)
}

func updateUser(r *userRepository, t *testing.T) {
	err := r.UpdateUser(&User{1, "sujeto1_act", "sujeto1@cambio"})
	if err != nil {
		t.Errorf("Error actualizando usuario 1: %v", err)
	}

	err = r.UpdateUser(&User{2, "sujeto2_act", "sujeto2@cambio"})
	if err != nil {
		t.Errorf("Error actualizando usuario 2: %v", err)
	}

	err = r.UpdateUser(&User{3, "sujeto3_act", "sujeto3@cambio"})
	if err != nil {
		t.Errorf("Error actualizando usuario 3: %v", err)
	}
}

func listUsers(r *userRepository, t *testing.T) {
	_, err := r.ListUsers()
	if err != nil {
		t.Errorf("Error listando usuarios: %v", err)
	}
}

func deleteUser(r *userRepository, t *testing.T) {
	err := r.DeleteUser(1)
	if err != nil {
		t.Errorf("Error eliminando usuario 1: %v", err)
	}
	_, err = r.GetUserByID(1)
	if err != sql.ErrNoRows {
		fmt.Println("El usuario 1 no se elimino correctamente")
	}

	err = r.DeleteUser(2)
	if err != nil {
		t.Errorf("Error eliminando usuario 2: %v", err)
	}
	_, err = r.GetUserByID(2)
	if err != sql.ErrNoRows {
		fmt.Println("El usuario 2 no se elimino correctamente")
	}

	err = r.DeleteUser(3)
	if err != nil {
		t.Errorf("Error eliminando usuario 3: %v", err)
	}
	_, err = r.GetUserByID(3)
	if err != sql.ErrNoRows {
		fmt.Println("El usuario 3 no se elimino correctamente")
	}
}
