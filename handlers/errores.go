package handlers

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

// Codigos SQLSTATE de Postgres
const (
	codigoViolacionFK     = "23503" // foreign_key_violation
	codigoViolacionUnique = "23505" // unique_violation
)

// Indica si err es un error de Postgres con el codigo SQLSTATE dado
func esErrorPg(err error, codigo string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == codigo
}

// Indica si err es una violacion de FK (ej: borrar un vendedor con articulos)
func esViolacionFK(err error) bool {
	return esErrorPg(err, codigoViolacionFK)
}

// Indica si err es una violacion de UNIQUE (ej: email duplicado)
func esViolacionUnique(err error) bool {
	return esErrorPg(err, codigoViolacionUnique)
}
