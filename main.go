package main

import (
    "fmt"
    "net/http"
)

func main() {
    staticDir := "./static"
    fileServer := http.FileServer(http.Dir(staticDir))
    http.Handle("/", fileServer)
    port := ":8080"

    fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

    err := http.ListenAndServe(port, nil)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
}

// func serve(w http.ResponseWriter, r *http.Request) {
//     if r.URL.Path != "/" || r.Method != http.MethodGet {
//         http.NotFound(w, r)
//         return
//     }
//     w.Header().Set("Content-Type", "text/html; charset=utf-8")
//     fmt.Fprint(w, `<!DOCTYPE html>
// <html>
// <head><title>Inicio</title></head>
// <body>
//   <h1>Bienvenido a ...</h1>
//   <p>Esta es la pagina principal.</p>
//   <!-- <a href=""><> -->
// </body>
// </html>
// `)
// }
