package logic

import (
	"net/http"
	"strconv"
)

func IsHtmxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

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

func writeAuthRequiredFragment(w http.ResponseWriter, loginURL string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := `<div class="bloque pixel-frame">
  <h1 class="display">ACCESO REQUERIDO</h1>
  <p>Para acceder a esta página tenés que iniciar sesión.</p>
  <a class="btn" href=` + loginURL + `>Iniciar sesión ↗</a>
</div>`
	w.Write([]byte(html))
}

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

// CarritoCookie devuelve las unidades del carrito guardadas en la cookie.
// El header la usa para elegir el ícono sin consultar la base.
func CarritoCookie(r *http.Request) int {
	c, err := r.Cookie("cart")
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(c.Value)
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

// SetCarritoCookie guarda la cantidad del carrito; 0 la invalida.
func SetCarritoCookie(w http.ResponseWriter, cantidad int) {
	valor := "0"
	maxAge := -1
	if cantidad > 0 {
		valor = strconv.Itoa(cantidad)
		maxAge = 0
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "cart",
		Value:    valor,
		Path:     "/",
		MaxAge:   maxAge,
		SameSite: http.SameSiteLaxMode,
	})
}
