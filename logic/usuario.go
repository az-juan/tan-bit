package logic

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
)

// ValidarEmail verifica el formato del email.
func ValidarEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

// ActualizarUsuarioInput son los campos editables de la cuenta.
type ActualizarUsuarioInput struct {
	Nombre   string
	Apellido string
	Email    string
	Password string // vacío = no cambiar
}

// ParseActualizarUsuario valida el formulario de datos personales.
func ParseActualizarUsuario(r *http.Request) (ActualizarUsuarioInput, bool) {
	in := ActualizarUsuarioInput{
		Nombre:   strings.TrimSpace(r.FormValue("nombre")),
		Apellido: strings.TrimSpace(r.FormValue("apellido")),
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: r.FormValue("password_nueva"),
	}
	if in.Nombre == "" || in.Apellido == "" || !ValidarEmail(in.Email) {
		return in, false
	}
	if in.Password != "" && len(in.Password) < 6 {
		return in, false
	}
	return in, true
}

// ActualizarVendedorInput son los campos editables del perfil vendedor.
type ActualizarVendedorInput struct {
	Telefono     int64
	Cbu          string
	Alias        string
	DireccionXmr sql.NullString
}

// ParseActualizarVendedor valida el formulario de datos complementarios.
func ParseActualizarVendedor(r *http.Request) (ActualizarVendedorInput, bool) {
	in := ActualizarVendedorInput{}
	telefono := digitsOnly(r.FormValue("telefono"))
	if telefono == "" {
		return in, false
	}
	t, err := strconv.ParseInt(telefono, 10, 64)
	if err != nil || t <= 0 {
		return in, false
	}
	cbu := strings.ReplaceAll(r.FormValue("cbu"), " ", "")
	cbu = strings.ReplaceAll(cbu, "-", "")
	if !ValidateCBU(cbu) {
		return in, false
	}
	alias := strings.TrimSpace(r.FormValue("alias"))
	if alias == "" {
		return in, false
	}
	dir := strings.TrimSpace(r.FormValue("direccion_xmr"))
	in.Telefono = t
	in.Cbu = cbu
	in.Alias = alias
	in.DireccionXmr = sql.NullString{String: dir, Valid: dir != ""}
	return in, true
}
