package main

import (
	"net/http"
	"net/http/httptest" //Paquete std para probar web servers sin abrir puertos reales
	"testing"           //Obligatorio en _test.go
)

func TestServidorEstatico(t *testing.T) { //Firma std en _test.go (Test+Mayus y recibe un puntero)
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))

	// Petición HTTP simulada a la raíz "/"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Error al crear la petición: %v", err) //Fatal para arranque
	}

	// Creamos un grabador de respuesta q captura lo que devuelve el handler
	rr := httptest.NewRecorder()

	fileServer.ServeHTTP(rr, req)

	// Verificamos que el código de estado HTTP sea 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Código de estado incorrecto: obtenido %v, esperado %v", status, http.StatusOK)
	}
}
