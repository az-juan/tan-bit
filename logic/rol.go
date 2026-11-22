package logic

import "net/http"

// RolCookie devuelve el rol guardado en la cookie: "vendedor", "usuario" o "".
// El header lo usa para elegir el ícono de cuenta sin consultar la base.
func RolCookie(r *http.Request) string {
	c, err := r.Cookie("role")
	if err != nil {
		return ""
	}
	switch c.Value {
	case "vendedor", "usuario":
		return c.Value
	default:
		return ""
	}
}

// SetRolCookie guarda el rol del usuario.
func SetRolCookie(w http.ResponseWriter, rol string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "role",
		Value:    rol,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearRolCookie invalida la cookie de rol.
func ClearRolCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "role",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
}
