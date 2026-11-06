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
	// if cat := r.URL.Query().Get("cat"); cat != "" {
	// 	filtered := arts[:0]
	// 	for _, a := range arts {
	// 		if string(a.Categoria) == cat {
	// 			filtered = append(filtered, a)
	// 		}
	// 	}
	// 	arts = filtered
	// }
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Marketplace(arts).Render(r.Context(), w)
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
	views.Favoritos().Render(r.Context(), w)
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
	views.Vender().Render(r.Context(), w)
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

func (h *ArticuloHandler) CreadoresHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Creadores().Render(r.Context(), w)
}
