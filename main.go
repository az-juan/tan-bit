package main

import (
	"log"
	"fmt"
	"net/http"
	"tan-bit.com/tan-bit/handlers"
	"tan-bit.com/tan-bit/db/connection"
	sqlc "tan-bit.com/tan-bit/db/sqlc"
)

func main() {
	db, err := connection.Open()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db) //Instanciamos queries para ejecutarse en nuestro pool
	artHandler := handlers.NewArticuloHandler(queries)
	vendHandler := handlers.NewVendedorHandler(queries)
	userHandler := handlers.NewUsuarioHandler(queries)


	staticDir := "./static"                            //Maneja automaticamente los Content-Type
	fileServer := http.FileServer(http.Dir(staticDir)) //de los archivos que sirve

	// Configurar rutas
	http.HandleFunc("/articulos", artHandler.ArtsHandler)
	http.HandleFunc("/articulos/", artHandler.ArtHandler)
	http.HandleFunc("/vendedores", vendHandler.VendsHandler)
    http.HandleFunc("/usuarios", userHandler.UsersHandler)
	http.HandleFunc("/health", healthHandler)

	http.Handle("GET /{$}", fileServer)
	http.HandleFunc("/", notFoundHandler)
	port := ":8080"

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Servidor andando!"))
	default:
		http.Error(w, "Metodo no permitido",
			http.StatusMethodNotAllowed)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "La ruta %s es inválida, por favor, vuelve a la anterior.\n", r.URL.Path)
}
