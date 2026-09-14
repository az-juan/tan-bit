package main

import (
	"context"         //Gestiona tiempos de vida y cancelaciones de peticiones a la base de datos
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"              //Lectura de variables de entorno del sistema operativo

	_ "github.com/jackc/pgx/v5/stdlib"  //driver estándar para postgreSQL
	db "tp1.com/aplicacion_web/db/sqlc"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5433"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "ayudanoentiendo"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "tan_bit"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				dbHost, dbPort, dbUser, dbPass, dbName)

	dbConn, err := sql.Open("pgx", connStr)     //Inicializa el pool de conexiones usando pgx
	if err != nil {
		log.Fatalf("Error al configurar conexion a la BD: %v", err)
	}
	defer dbConn.Close()                        //Cierre del pool de conexiones

	if err := dbConn.Ping(); err != nil {
		log.Printf("Aviso: No se pudo contactar la BD al inicio: %v", err)
	} else {
		fmt.Println("Conexión con PostgreSQL establecida correctamente.")
	}

	queries := db.New(dbConn)                   //Instancia el cliente generado por sqlc

	http.HandleFunc("/api/productos", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		productos, err := queries.ListProducts(ctx, db.ListProductsParams{
			Limit:  20,
			Offset: 0,
		})
		if err != nil {
			http.Error(w, "Error al consultar productos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(productos)
	})

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)

	port := ":8080"
	fmt.Printf("Servidor TAN-BIT escuchando en http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error al iniciar el servidor web: %v", err)
	}
}
