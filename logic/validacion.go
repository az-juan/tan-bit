// Package logic contiene funciones de parseo, validación y autenticación
// comunes usadas en los handlers.
package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
)

// ParseFormSafe parsea el form de la request, devuelve error si falla.
func ParseFormSafe(r *http.Request) error {
	return r.ParseForm()
}

// GetFormValue obtiene un valor del form, o retorna string vacío si no existe.
func GetFormValue(r *http.Request, key string) string {
	return r.FormValue(key)
}

// IsHtmxRequest retorna true si la petición proviene de htmx.
func IsHtmxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// GetCookieOrRedirect retorna el cookie usuario_id, o redirige a /login si no existe.
// Retorna (id, ok) donde ok es false si no hay cookie válida.
//
// Cuando la petición es de htmx y el usuario no está autenticado, en vez de redirigir
// (que forzaría una carga completa de la página y repetir header/footer) escribe un
// fragmento HTML que htmx puede intercambiar por el elemento #main-content, y
// retorna (0, false). El handler que llamó debe retornar inmediatamente tras este
// valor, sin continuar renderizando la página.
func GetCookieOrRedirect(w http.ResponseWriter, r *http.Request) (int, bool) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		if IsHtmxRequest(r) {
			writeAuthRequiredFragment(w, "/login")
			return 0, false
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 0, false
	}
	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		if IsHtmxRequest(r) {
			writeAuthRequiredFragment(w, "/login")
			return 0, false
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 0, false
	}
	return id, true
}

// writeAuthRequiredFragment escribe un fragmento HTML para htmx que se intercambia
// en el elemento #main-content cuando el usuario no está autenticado.
// Debe ser el único usuario de w; el handler que lo llamó debe retornar.
func writeAuthRequiredFragment(w http.ResponseWriter, loginURL string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := `<div class="auth-required pixel-frame">
  <h1 class="display">ACCESO REQUERIDO</h1>
  <p>Para acceder a esta página tenés que iniciar sesión.</p>
  <a class="btn" href=` + loginURL + `>Iniciar sesión ↗</a>
</div>`
	w.Write([]byte(html))
}

// CheckAuth verifica si el usuario está autenticado.
// Retorna (ok bool, id int).
func CheckAuth(r *http.Request) (bool, int) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return false, 0
	}
	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return false, 0
	}
	return true, id
}

// ValidateEmail verifica que el email tenga formato válido.
func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ValidatePassword verifica que la contraseña tenga al menos 6 caracteres.
func ValidatePassword(password string) bool {
	return len(password) >= 6
}

// ValidateName verifica que el nombre no esté vacío.
func ValidateName(name string) bool {
	return strings.TrimSpace(name) != ""
}

// ValidateRequiredFields verifica que todos los campos obligatorios estén presentes.
func ValidateRequiredFields(fields map[string]string) bool {
	for _, v := range fields {
		if strings.TrimSpace(v) == "" {
			return false
		}
	}
	return true
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

// ValidateRole verifica que el rol sea válido ("usuario" o "vendedor").
func ValidateRole(role string) bool {
	return role == "usuario" || role == "vendedor"
}

// ValidateTerminos verifica que los términos hayan sido aceptados.
func ValidateTerminos(terminos string) bool {
	return terminos == "on"
}

// HashPassword crea un hash simple de la contraseña (reemplazar con bcrypt en producción).
func HashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

// IsValidString verifica que un string no esté vacío después de trim.
func IsValidString(s string) bool {
	return strings.TrimSpace(s) != ""
}
