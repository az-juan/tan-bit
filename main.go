package main

import (
	"fmt"
	"log"
	"net/http"
	"tan-bit.com/tan-bit/db/connection"
	sqlc "tan-bit.com/tan-bit/db/sqlc"
	"tan-bit.com/tan-bit/handlers"
)

func main() {
	db, err := connection.Open()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db) //Instanciamos queries para ejecutarse en nuestro pool
	artHandler := handlers.NewArticuloHandler(queries)
	userHandler := handlers.NewUsuarioHandler(queries)
	vendHandler := handlers.NewVendedorHandler(queries)

	fs := http.FileServer(http.Dir("static"))
	us := http.FileServer(http.Dir("db/uploads"))

	// Configurar rutas
	http.HandleFunc("/health", handlers.HealthHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", us))
	http.Handle("/vender", vendHandler)
	http.Handle("/login", userHandler)
	http.Handle("/registrarse", userHandler)
	http.Handle("/mi_cuenta", userHandler)
	http.Handle("/", artHandler)
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
