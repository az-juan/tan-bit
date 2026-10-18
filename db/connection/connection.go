package connection

import (
	"fmt"
	"charm.land/log/v2"
	"os"
	"database/sql" //Estandar para interactuar con la db
	_ "github.com/jackc/pgx/v5/stdlib" //Driver estándar para postgreSQL
	"github.com/joho/godotenv"
)

func Open() (*sql.DB, error) {
	if err := godotenv.Load("../example.env"); err != nil {
		log.Fatalf("Archivo .env no encontrado")
	}
	connStr := fmt.Sprintf( //Leemos el .env
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("pgx", connStr) //Inicializa pool de conexiones a la db usando "pgx"
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}