// Package logic contiene funciones de parseo, validación y autenticación
// comunes usadas en los handlers.
package logic

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// CodigoViolacionFK es el código SQLSTATE de la violación de foreign key.
const CodigoViolacionFK = "23503"

// CodigoViolacionUnique es el código SQLSTATE de la violación de unique.
const CodigoViolacionUnique = "23505"

// EsErrorPg indica si err es un error de Postgres con el código SQLSTATE dado.
func EsErrorPg(err error, codigo string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == codigo
}

// EsViolacionFK indica si err es una violación de FK.
func EsViolacionFK(err error) bool {
	return EsErrorPg(err, CodigoViolacionFK)
}

// EsViolacionUnique indica si err es una violación de UNIQUE.
func EsViolacionUnique(err error) bool {
	return EsErrorPg(err, CodigoViolacionUnique)
}

// ColumnaViolacionUnique detecta qué columna generó la violación de UNIQUE.
func ColumnaViolacionUnique(err error) string {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return ""
	}
	switch {
	case strings.Contains(strings.ToLower(pgErr.ConstraintName), "usuario_email") || strings.Contains(strings.ToLower(pgErr.ConstraintName), "email"):
		return "email"
	case strings.Contains(strings.ToLower(pgErr.ConstraintName), "vendedor_cbu"):
		return "cbu"
	case strings.Contains(strings.ToLower(pgErr.ConstraintName), "vendedor_alias"):
		return "alias"
	case strings.Contains(strings.ToLower(pgErr.ConstraintName), "vendedor_telefono"):
		return "telefono"
	}
	return ""
}
