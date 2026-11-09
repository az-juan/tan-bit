package handlers

import (
	"fmt"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Servidor andando!"))
	default:
		http.Error(w, "Metodo no permitido",
			http.StatusMethodNotAllowed)
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "La ruta %s es inválida, por favor, vuelve a la anterior.\n", r.URL.Path)
}
