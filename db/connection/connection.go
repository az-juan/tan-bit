package connection

import (
	"database/sql" //Estandar para interactuar con la db
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib" //Driver estándar para postgreSQL
	"os"
)

func Open() (*sql.DB, error) {
	connStr := fmt.Sprintf( //Leemos el .env
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
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
