// Package logic contiene funciones de parseo, validación y autenticación
// comunes usadas en los handlers.
package logic

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
)

// RegistrarseInput representa los datos parseados y validados del formulario
// de registro de usuario/vendedor.
type RegistrarseInput struct {
	Email        string
	Password     string
	Nombre       string
	Apellido     string
	Role         string // normalizado: "usuario" o "vendedor"
	Telefono     string // solo dígitos numéricos (vendedor)
	CBU          string // solo dígitos numéricos (13 caracteres)
	Alias        string
	DireccionXmr string
	Terminos     string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

func ParseRegistrarse(r *http.Request) *RegistrarseInput {
	email := strings.TrimSpace(r.FormValue("email"))
	if !emailRegex.MatchString(email) {
		return nil
	}

	password := r.FormValue("password")
	if len(password) < 6 {
		return nil
	}
	if r.FormValue("password2") != password {
		return nil
	}

	nombre := strings.TrimSpace(r.FormValue("nombre"))
	if nombre == "" {
		return nil
	}
	apellido := strings.TrimSpace(r.FormValue("apellido"))
	if apellido == "" {
		return nil
	}

	role := r.FormValue("role")
	if role != "vendedor" && role != "usuario" {
		role = "usuario"
	}

	telefono := r.FormValue("telefono")
	cbu := r.FormValue("cbu")
	alias := strings.TrimSpace(r.FormValue("alias"))
	direccionXmr := r.FormValue("direccion_xmr")

	if role == "vendedor" {
		if !ValidatePhoneNumber(telefono) {
			return nil
		}
		telefono = digitsOnly(telefono)
		if telefono == "" {
			return nil
		}
		cbu = strings.ReplaceAll(cbu, " ", "")
		cbu = strings.ReplaceAll(cbu, "-", "")
		if !ValidateCBU(cbu) {
			return nil
		}
		if alias == "" {
			return nil
		}
	}

	terminos := r.FormValue("terminos")
	if terminos != "on" {
		return nil
	}

	return &RegistrarseInput{
		Email:        email,
		Password:     password,
		Nombre:       nombre,
		Apellido:     apellido,
		Role:         role,
		Telefono:     telefono,
		CBU:          cbu,
		Alias:        alias,
		DireccionXmr: direccionXmr,
		Terminos:     terminos,
	}
}

// digitsOnly devuelve solo los caracteres dígitos de s.
func digitsOnly(s string) string {
	var b strings.Builder
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// ValidatePhoneNumber verifica que el teléfono contenga solo dígitos, +, espacios y guiones.
func ValidatePhoneNumber(phone string) bool {
	for _, ch := range phone {
		if !((ch >= '0' && ch <= '9') || ch == '+' || ch == ' ' || ch == '-') {
			return false
		}
	}
	return true
}

// ValidateCBU verifica que el CBU tenga 13 dígitos numéricos.
func ValidateCBU(cbu string) bool {
	cbu = strings.ReplaceAll(cbu, " ", "")
	cbu = strings.ReplaceAll(cbu, "-", "")
	if len(cbu) != 13 {
		return false
	}
	for _, ch := range cbu {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// HashPassword crea un hash bcrypt de la contraseña.
func HashPassword(password string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(h)
}

// CheckPassword verifica la contraseña contra el hash bcrypt guardado.
func CheckPassword(stored string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(stored), []byte(password)) == nil
}

// IsValidString verifica que un string no esté vacío después de trim.
func IsValidString(s string) bool {
	return strings.TrimSpace(s) != ""
}
