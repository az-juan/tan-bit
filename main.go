package main

import (
    "fmt"
    "net/http"
)

func main() {
    staticDir := "./static"
    fileServer := http.FileServer(http.Dir(staticDir))
    http.Handle("/", fileServer)  //Maneja automaticamente los Content-Type
    port := ":8080"               //de los archivos que sirve

    fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

    err := http.ListenAndServe(port, nil)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
}