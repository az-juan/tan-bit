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
	
	http.Handle("GET /{$}", fileServer)   	//Maneja automaticamente los Content-Type
	   										//de los archivos que sirve
	http.HandleFunc("/", notFoundHandler)
	port := ":8080"              

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
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
    fmt.Fprintf(w, "La ruta %s es inválida, por favor, vuelve a la anterior.%s\n", r.URL.Path)
}