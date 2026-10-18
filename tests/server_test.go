package tests

import (
	"net/http"
	"testing"           //Obligatorio en _test.go
)

func TestAPI(t *testing.T) { //Firma std en _test.go (Test+Mayus y recibe un puntero)
	// Petición HTTP simulada a la raíz "/"
	req, err := http.NewRequest("GET", "http://localhost:8080/health", nil)
	if err != nil {
		t.Fatalf("Error al crear la petición: %v", err) //Fatal para arranque
	}
	
	// Abre socket TCP para solicitar, luego lo cierra
	ans, err := http.DefaultClient.Do(req)
        if err != nil {
            t.Fatalf("Error al conectar con el servidor en el contenedor: %v", err)
        }
	defer ans.Body.Close()

	// Verificamos que el código de estado HTTP sea 200 OK
	if status := ans.StatusCode; status != http.StatusOK {
		t.Errorf("Código de estado incorrecto: obtenido %v, esperado %v", status, http.StatusOK)
	}
}
