package main

import (
	"fmt"
	_ "github.com/a-h/templ"
	"net/http"
	"tan-bit.com/tan-bit/handlers"
)

func main() {
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))

	// Configurar rutas
	http.HandleFunc("/articulos", handlers.ArtsHandler)
	http.HandleFunc("/articulos/", handlers.ArtHandler)
	http.HandleFunc("/health", healthHandler)
	http.Handle("/", fileServer) //Maneja automaticamente los Content-Type
	port := ":8080"              //de los archivos que sirve

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Servidor andando!"))
}
