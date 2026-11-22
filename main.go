package main

import (
	"fmt"
	"log"
	"net/http"
	"tan-bit.com/tan-bit/db/connection"
	"tan-bit.com/tan-bit/handlers"
	"tan-bit.com/tan-bit/logic"
)

func main() {
	logic.StartNotificationDispatcher()
	db, err := connection.Open()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	datos := logic.NuevoDatos(db) // acceso a la base vía sqlc
	artHandler := handlers.NewArticuloHandler(datos)
	userHandler := handlers.NewUsuarioHandler(datos)
	vendHandler := handlers.NewVendedorHandler(datos)

	fs := http.FileServer(http.Dir("static"))
	us := http.FileServer(http.Dir("uploads"))

	// Configurar rutas
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/evento", handlers.EventoHandler)
	http.HandleFunc("/privacidad", handlers.PrivacidadHandler)
	http.HandleFunc("/terminos", handlers.TerminosHandler)
	http.HandleFunc("/historia", handlers.HistoriaHandler)
	http.HandleFunc("/creadores", handlers.CreadoresHandler)
	http.HandleFunc("/contacto", handlers.ContactoHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", us))
	http.Handle("/vender", vendHandler)
	http.Handle("/registrarse", userHandler)
	http.Handle("/login", userHandler)
	http.Handle("/logout", userHandler)
	http.Handle("/mi_cuenta", userHandler)
	http.Handle("/mi_cuenta/delete", userHandler)
	http.Handle("/search", artHandler)
	http.Handle("/", artHandler)
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
