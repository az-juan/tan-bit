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
	html := `<div class="auth-required pixel-frame">
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
