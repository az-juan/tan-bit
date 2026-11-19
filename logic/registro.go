// Package logic contiene funciones de parseo, validación y autenticación
// comunes usadas en los handlers.
package logic

import (
	"net/http"
	"strings"
)

// RegistrarseInput representa los datos parseados y validados del formulario
// de registro de usuario/vendedor.
type RegistrarseInput struct {
	Email    string
	Password string
	Nombre   string
	Apellido string
	Role     string // normalizado: "usuario" o "vendedor"
	Telefono string // solo dígitos numéricos (vendedor)
	CBU      string // solo dígitos numéricos (13 caracteres)
	Alias    string
	Terminos string
}

// ParseRegistrarse parsea el formulario de registro y lo valida.
// Devuelve el dato parseado y un mensaje de error (vacío si es válido).
// La lógica de validación (if) se mueve aquí, fuera del handler.
func ParseRegistrarse(r *http.Request) (*RegistrarseInput, string) {
	email := GetFormValue(r, "email")
	password := GetFormValue(r, "password")
	nombre := GetFormValue(r, "nombre")
	apellido := GetFormValue(r, "apellido")
	role := GetFormValue(r, "role")
	if !ValidateRole(role) {
		role = "usuario"
	}
	telefono := GetFormValue(r, "telefono")
	cbu := GetFormValue(r, "cbu")
	alias := GetFormValue(r, "alias")
	terminos := GetFormValue(r, "terminos")

	input := &RegistrarseInput{
		Email:    email,
		Password: password,
		Nombre:   nombre,
		Apellido: apellido,
		Role:     role,
		Telefono: telefono,
		CBU:      cbu,
		Alias:    alias,
		Terminos: terminos,
	}

	if !ValidateRequiredFields(map[string]string{
		"nombre":   input.Nombre,
		"apellido": input.Apellido,
		"email":    input.Email,
		"password": input.Password,
	}) {
		return input, "COMPLETÁ TODOS LOS CAMPOS OBLIGATORIOS."
	}
	if !ValidateEmail(input.Email) {
		return input, "EMAIL no válido."
	}
	if !ValidatePassword(input.Password) {
		return input, "LA CONTRASEÑA debe tener al menos 6 caracteres."
	}
	if !ValidateTerminos(input.Terminos) {
		return input, "Debés aceptar los términos del servicio y la política de privacidad."
	}

	// Validación de datos complementarios de vendedor.
	if input.Role == "vendedor" {
		if !ValidateRequiredFields(map[string]string{
			"telefono": input.Telefono,
			"cbu":      input.CBU,
			"alias":    input.Alias,
		}) {
			return input, "Para registrarte como vendedor necesitás completar TELÉFONO, CBU y ALIAS."
		}
		if !ValidatePhoneNumber(input.Telefono) {
			return input, "TELÉFONO solo puede contener dígitos, + y espacios."
		}
		cbuClean := strings.ReplaceAll(input.CBU, " ", "")
		cbuClean = strings.ReplaceAll(cbuClean, "-", "")
		if !ValidateCBU(cbuClean) {
			return input, "EL CBU debe tener 13 dígitos numéricos."
		}
		input.CBU = cbuClean
		input.Telefono = digitsOnly(input.Telefono)
	}

	return input, ""
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
