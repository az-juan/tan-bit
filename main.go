package main

import (
	"context" //Gestiona tiempos de vida y cancelaciones de peticiones a la base de datos
	"database/sql" //Estandar para interactuar con la db
	"log"
	"fmt"
	"net/http"
	
	sqlc "aplicacion_web/db/sqlc"      //Main sabe a donde ir a buscar el sqlc
	_ "github.com/jackc/pgx/v5/stdlib" //Driver estándar para postgreSQL
)

func main() {
	connStr := "user=postgres password=ayudanoentiendo dbname=tan_bit"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	defer db.Close()

	queries := sqlc.New(db)
	ctx := context.Background()

	createdUser, err := queries.CreateUser(ctx, // Create
		sqlc.CreateUserParams{
			Nombre:   "Nikola",
			Apellido: "Tesla",
			Email:    "torreTesla@example.com",
		})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	fmt.Printf("Created user: %+v\n", createdUser)

	user, err := queries.GetUserByID(ctx, createdUser.IDUsuario) // Read One
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	fmt.Printf("Retrieved user: %+v\n", user)

	users, err := queries.ListUsers(ctx) // Read Many
	if err != nil {
		log.Fatalf("failed to list users: %v", err)
	}
	fmt.Printf("All users: %+v\n", users)

	err = queries.UpdateUser(ctx, sqlc.UpdateUserParams{ // Update
		IDUsuario: createdUser.IDUsuario,
		Nombre:     "Albert",
		Apellido: "Einstein",
		Email:    "montapuercos@example.com",
	})
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}
	fmt.Println("User updated successfully")

	updatedUser, err := queries.GetUserByID(ctx, createdUser.IDUsuario)
	if err != nil {
		log.Fatalf("failed to get updated user: %v", err)
	}
	fmt.Printf("Updated user: %+v\n", updatedUser)

	err = queries.DeleteUser(ctx, createdUser.IDUsuario) // Delete
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}
	fmt.Println("User deleted successfully")

	_, err = queries.GetUserByID(ctx, createdUser.IDUsuario)
	if err == sql.ErrNoRows {
		fmt.Println("User not found after deletion")
	} else if err != nil {
		log.Fatalf("failed to get user after deletion: %v", err)
	}

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)

	port := ":8080"
	fmt.Printf("Servidor TAN-BIT escuchando en http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error al iniciar el servidor web: %v", err)
	}
}
