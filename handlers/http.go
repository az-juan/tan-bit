package handlers

import (
	"fmt"
	"net/http"
	"tan-bit.com/tan-bit/views"
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

func EventoHandler(w http.ResponseWriter, r *http.Request) {
	views.Evento().Render(r.Context(), w)
}

func PrivacidadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Privacidad().Render(r.Context(), w)
}

func ContactoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Contacto().Render(r.Context(), w)
}

func TerminosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Terminos().Render(r.Context(), w)
}

func HistoriaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Historia().Render(r.Context(), w)
}

func CreadoresHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Creadores().Render(r.Context(), w)
}
