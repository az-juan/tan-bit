package handlers

import (
	"database/sql"
	"net/http"
	db "tan-bit.com/tan-bit/db/sqlc"
	"tan-bit.com/tan-bit/logic"
	"tan-bit.com/tan-bit/views"
)

type VendedorHandler struct {
	datos *logic.Datos
}

func NewVendedorHandler(datos *logic.Datos) *VendedorHandler {
	return &VendedorHandler{datos: datos}
}

// ServeHTTP para vendedor handler
func (h *VendedorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/vender":
		h.VenderHandler(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

// Manejador de página de venta
func (h *VendedorHandler) VenderHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	v, err := h.datos.Queries().GetVendedorByID(r.Context(), int32(id))
	if err != nil || v.Telefono == 0 {
		notificar(w, r, logic.NotificationWarning, "atención", "Necesitás una cuenta de vendedor", "Completá tus datos de vendedor para publicar piezas.", "/registrarse", "Ir a registrarme")
		http.Redirect(w, r, "/registrarse", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			r.ParseForm()
		}
		data, err := logic.ParseVender(r, id)
		if err != nil {
			notificar(w, r, logic.NotificationError, "error", "No pudimos publicar", "Revisá los campos obligatorios: título, descripción, ubicación, precio, stock, tipo y zona de entrega.", "/vender", "Reintentar")
			http.Redirect(w, r, "/vender", http.StatusSeeOther)
			return
		}
		params := db.CreateArticuloParams{
			Titulo:      data.Titulo,
			Descripcion: data.Descripcion,
			Condicion:   db.Condicion(data.Condicion),
			Categoria:   db.Categoria(data.Categoria),
			Precio:      data.Precio,
			PrecioXmr:   data.PrecioXmr,
			Stock:       data.Stock,
			IDVendedor:  int32(id),
			Ubicacion:   data.Ubicacion,
			TipoEntrega: db.Entrega(data.TipoEntrega),
			CostoEnvio:  sql.NullFloat64{Float64: data.CostoEnvio, Valid: data.CostoEnvio >= 0},
			ZonaEntrega: sql.NullString{String: data.ZonaEntrega, Valid: data.ZonaEntrega != ""},
			RutaImagen:  sql.NullString{String: data.RutaImagen, Valid: data.RutaImagen != ""},
		}
		if _, err := h.datos.Queries().CreateArticulo(r.Context(), params); err != nil {
			notificar(w, r, logic.NotificationError, "error", "No pudimos publicar", "Ocurrió un error al guardar la publicación.", "/vender", "Reintentar")
			http.Redirect(w, r, "/vender", http.StatusSeeOther)
			return
		}
		notificar(w, r, logic.NotificationSuccess, "éxito", "Artículo publicado", data.Titulo, "/", "Ver publicación")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	views.Pagina(w, r, views.Vender(logic.VenderInput{}, logic.CategoriasDB(), logic.CondicionesDB())).Render(r.Context(), w)
}
