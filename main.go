package main

import (
	"fmt"
	"net/http"
)

func main() {
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir)) //Creamos un handler de archivos html

	http.Handle("/", fileServer) //Lo registramos para que atienda peticiones a /
	//Maneja automaticamente los Content-Type
	//de los archivos que sirve

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
