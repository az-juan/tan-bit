package handlers

import (
	"net/http"
	"tan-bit.com/tan-bit/views"
)

func (h *ArticuloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.MarketplaceHandler(w, r)
	case "/vender":
		h.VenderHandler(w, r)
	case "/carrito":
		h.CarritoHandler(w, r)
	case "/favoritos":
		h.FavoritosHandler(w, r)
	case "/evento":
		h.EventoHandler(w, r)
	case "/login":
		h.LoginHandler(w, r)
	case "/privacidad":
		h.PrivacidadHandler(w, r)
	case "/marketplace":
		h.MarketplaceHandler(w, r)
	case "/contacto":
		h.ContactoHandler(w, r)
	case "/terminos":
		h.TerminosHandler(w, r)
	case "/historia":
		h.HistoriaHandler(w, r)
	case "/precio-entrega":
		h.PrecioEntregaHandler(w, r)
	case "/revisar-publicar":
		h.RevisarPublicarHandler(w, r)
	case "/articulos":
		h.ArtsHandler(w, r)
	case "/registrarse":
		h.RegistrarseHandler(w, r)
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
func (h *ArticuloHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Login().Render(r.Context(), w)
}
func (h *ArticuloHandler) PrivacidadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Privacidad().Render(r.Context(), w)
}
func (h *ArticuloHandler) VenderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Vender(views.VenderData{}, views.CategoriasDB(), views.CondicionesDB()).Render(r.Context(), w)
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

func (h *ArticuloHandler) PrecioEntregaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulario inválido", http.StatusBadRequest)
			return
		}
		data := views.ParseVenderData(r)
		http.Redirect(w, r, "/precio-entrega?"+data.Query(), http.StatusSeeOther)
		return
	}
	data := views.ParseVenderData(r)
	views.PrecioEntrega(data).Render(r.Context(), w)
}

func (h *ArticuloHandler) RevisarPublicarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulario inválido", http.StatusBadRequest)
			return
		}
		data := views.ParseVenderData(r)
		http.Redirect(w, r, "/revisar-publicar?"+data.Query(), http.StatusSeeOther)
		return
	}
	data := views.ParseVenderData(r)
	views.RevisarPublicar(data).Render(r.Context(), w)
}

func (h *ArticuloHandler) RegistrarseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Registrarse().Render(r.Context(), w)
}
func (h *ArticuloHandler) CreadoresHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Creadores().Render(r.Context(), w)
}
