package main

import (
	"charm.land/log/v2"
	"database/sql"
	"fmt"
	_ "github.com/a-h/templ"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"net/http"
	"os"
	"tan-bit.com/tan-bit/internal/middleware"
	"tan-bit.com/tan-bit/pkg/handlers"
)

var nc *nats.Conn

func main() {
	execDB()
	handlers.ExecBroker()
	go middleware.ExecConsumer()

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

func execDB() {
	if err := godotenv.Load("./example.env"); err != nil {
		log.Fatalf("Archivo .env no encontrado")
	}
	connStr := "host=database port=5432 user=" + os.Getenv("PG_USER") + " password=" + os.Getenv("PG_PASSWD") + " database=" + os.Getenv("DB_NAME")
	db, err := sql.Open("pgx", connStr)
	db.Ping()
	if err != nil {
		log.Fatalf("Error al conectarse a la base de datos: %v", err)
	}
	defer db.Close()

}
