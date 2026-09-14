package main

import (
	"database/sql"
	"fmt"
	_ "github.com/a-h/templ"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load("example.env"); err != nil {
		log.Fatalf("no .env file found")
	}
	connStr := "host=database port=5432 user=" + os.Getenv("PG_USER") + " password=" + os.Getenv("PG_PASSWD") + " database=" + os.Getenv("PG_DB_NAME")
	db, err := sql.Open("pgx", connStr)
	db.Ping()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer) //Maneja automaticamente los Content-Type
	port := ":8080"              //de los archivos que sirve
	// http.HandleFunc("/*", handle404)

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
