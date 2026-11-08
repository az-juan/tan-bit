package handlers

import (
	"net/http"
	"strconv"

	db "tan-bit.com/tan-bit/db/sqlc"
	"tan-bit.com/tan-bit/views"
)

func (h *ArticuloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.MarketplaceHandler(w, r)
	case "/carrito":
		h.CarritoHandler(w, r)
	case "/favoritos":
		h.FavoritosHandler(w, r)
	case "/evento":
		h.EventoHandler(w, r)
	case "/privacidad":
		h.PrivacidadHandler(w, r)
	case "/contacto":
		h.ContactoHandler(w, r)
	case "/terminos":
		h.TerminosHandler(w, r)
	case "/historia":
		h.HistoriaHandler(w, r)
	case "/articulos":
		h.ArtsHandler(w, r)
	case "/articulos/":
		h.ArtsHandler(w, r)
	case "/registrarse":
		h.RegistrarseHandler(w, r)
	case "/mi_cuenta":
		h.MiCuentaHandler(w, r)
	case "/creadores":
		h.CreadoresHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *ArticuloHandler) MarketplaceHandler(w http.ResponseWriter, r *http.Request) {
	arts, err := h.queries.GetArticulos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
		return
	}
	views.Marketplace(arts, views.BuildFiltros(arts)).Render(r.Context(), w)
}

func (h *ArticuloHandler) CarritoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Carrito().Render(r.Context(), w)
}
func (h *ArticuloHandler) EventoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Evento().Render(r.Context(), w)
}
func (h *ArticuloHandler) FavoritosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	arts, err := h.queries.GetArticulos(r.Context())
	if err != nil {
		views.Favoritos(nil).Render(r.Context(), w)
		return
	}
	views.Favoritos(views.BuildFiltros(arts)).Render(r.Context(), w)
}
func (h *ArticuloHandler) PrivacidadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Privacidad().Render(r.Context(), w)
}

func (h *ArticuloHandler) ContactoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Contacto().Render(r.Context(), w)
}

func (h *ArticuloHandler) TerminosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Terminos().Render(r.Context(), w)
}

func (h *ArticuloHandler) HistoriaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Historia().Render(r.Context(), w)
}

func (h *ArticuloHandler) RegistrarseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.Registrarse().Render(r.Context(), w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}
	params := db.CreateUsuarioParams{
		Apellido: r.FormValue("apellido"),
		Nombre:   r.FormValue("nombre"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	nuevoUsuario, err := h.queries.CreateUsuario(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if tel := r.FormValue("telefono"); tel != "" {
		t, _ := strconv.ParseInt(tel, 10, 64)
		_, _ = h.queries.CreateVendedor(r.Context(), db.CreateVendedorParams{IDVendedor: nuevoUsuario.IDUsuario, Telefono: t})
	}
	http.SetCookie(w, &http.Cookie{Name: "user_id", Value: strconv.Itoa(int(nuevoUsuario.IDUsuario)), Path: "/"})
	http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
}
func (h *ArticuloHandler) MiCuentaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.MiCuenta().Render(r.Context(), w)
}
func (h *ArticuloHandler) CreadoresHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Creadores().Render(r.Context(), w)
}
