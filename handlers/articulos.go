package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	db "tan-bit.com/tan-bit/db/sqlc"
	"tan-bit.com/tan-bit/logic"
	"tan-bit.com/tan-bit/views"
)

type ArticuloHandler struct {
	datos *logic.Datos
}

func NewArticuloHandler(datos *logic.Datos) *ArticuloHandler {
	return &ArticuloHandler{datos: datos}
}

func (h *ArticuloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case path == "/":
		h.MarketplaceHandler(w, r)
	case path == "/carrito":
		h.CarritoHandler(w, r)
	case path == "/carrito/add":
		h.AddCarritoHandler(w, r)
	case strings.HasPrefix(path, "/carrito/"):
		h.RemoveCarritoHandler(w, r)
	case path == "/compra":
		h.CompraHandler(w, r)
	case path == "/favoritos":
		h.FavoritosHandler(w, r)
	case path == "/favoritos/add":
		h.AddFavoritoHandler(w, r)
	case strings.HasPrefix(path, "/favoritos/"):
		h.RemoveFavoritoHandler(w, r)
	case strings.HasPrefix(path, "/articulos/"):
		if strings.HasSuffix(path, "/editar") && r.Method == http.MethodGet {
			h.EditarArticuloHandler(w, r)
		} else if r.Method == http.MethodDelete {
			h.DeleteArticuloHandler(w, r)
		} else if r.Method == http.MethodPost {
			h.ActualizarArticuloHandler(w, r)
		} else if r.Method == http.MethodGet {
			h.FichaArticuloHandler(w, r)
		} else {
			http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		}
	case path == "/search":
		h.SearchHandler(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (h *ArticuloHandler) MarketplaceHandler(w http.ResponseWriter, r *http.Request) {
	arts, err := h.datos.Queries().GetArticulos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
		return
	}
	if categoria := r.FormValue("categoria"); categoria != "" {
		filtered := make([]db.GetArticulosRow, 0, len(arts))
		for _, a := range arts {
			if string(a.Categoria) == categoria {
				filtered = append(filtered, a)
			}
		}
		arts = filtered
	}
	data := make([]views.ArticuloData, len(arts))
	for i, a := range arts {
		data[i] = views.ArticuloData{
			IDArticulo:       a.IDArticulo,
			Titulo:           a.Titulo,
			Descripcion:      a.Descripcion,
			Condicion:        string(a.Condicion),
			Categoria:        string(a.Categoria),
			RutaImagen:       a.RutaImagen,
			Precio:           a.Precio,
			PrecioXmr:        a.PrecioXmr,
			Stock:            a.Stock,
			FechaPublicacion: a.FechaPublicacion,
			Ubicacion:        a.Ubicacion,
			TipoEntrega:      string(a.TipoEntrega),
			CostoEnvio:       a.CostoEnvio,
			ZonaEntrega:      a.ZonaEntrega,
			Alias:            "",
			Cbu:              "",
			DireccionXmr:     sql.NullString{},
			Cantidad:         0,
		}
	}
	if logic.IsHtmxRequest(r) {
		// Los filtros recargan solo la grilla, no la página entera.
		views.Pagina(w, r, views.GrillaMarketplace(data)).Render(r.Context(), w)
		return
	}
	views.Pagina(w, r, views.Marketplace(data, views.BuildFiltros(arts))).Render(r.Context(), w)
}

func (h *ArticuloHandler) CarritoHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	arts, err := h.datos.Queries().GetCarritoByUsuario(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
		return
	}
	var total float64
	var totalXmr float64
	data := make([]views.ArticuloData, len(arts))
	for i, a := range arts {
		total += a.Precio * float64(a.Cantidad)
		totalXmr += a.PrecioXmr * float64(a.Cantidad)
		data[i] = views.ArticuloData{
			IDArticulo:       a.IDArticulo,
			Titulo:           a.Titulo,
			Descripcion:      a.Descripcion,
			Condicion:        string(a.Condicion),
			Categoria:        string(a.Categoria),
			RutaImagen:       a.RutaImagen,
			Precio:           a.Precio,
			PrecioXmr:        a.PrecioXmr,
			Stock:            a.Stock,
			FechaPublicacion: a.FechaPublicacion,
			Ubicacion:        a.Ubicacion,
			TipoEntrega:      string(a.TipoEntrega),
			CostoEnvio:       a.CostoEnvio,
			ZonaEntrega:      a.ZonaEntrega,
			Alias:            "",
			Cbu:              "",
			DireccionXmr:     sql.NullString{},
			Cantidad:         a.Cantidad,
		}
	}
	views.Pagina(w, r, views.Carrito(data, total, totalXmr)).Render(r.Context(), w)
}

func (h *ArticuloHandler) CompraHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		arts, err := h.datos.Queries().GetCarritoCompra(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
			return
		}
		if len(arts) == 0 {
			http.Redirect(w, r, "/carrito", http.StatusSeeOther)
			return
		}
		entregas := make(map[string]bool)
		for _, a := range arts {
			entregas[string(a.TipoEntrega)] = true
		}
		entregaSlice := make([]string, 0, len(entregas))
		for k := range entregas {
			entregaSlice = append(entregaSlice, k)
		}
		sort.Strings(entregaSlice)
		var total float64
		var totalXmr float64
		var data = make([]views.ArticuloData, len(arts))
		for i, a := range arts {
			total += a.Precio * float64(a.Cantidad)
			totalXmr += a.PrecioXmr * float64(a.Cantidad)
			data[i] = views.ArticuloData{
				IDArticulo:       a.IDArticulo,
				Titulo:           a.Titulo,
				Descripcion:      a.Descripcion,
				Condicion:        string(a.Condicion),
				Categoria:        string(a.Categoria),
				RutaImagen:       a.RutaImagen,
				Precio:           a.Precio,
				PrecioXmr:        a.PrecioXmr,
				Stock:            a.Stock,
				FechaPublicacion: a.FechaPublicacion,
				Ubicacion:        a.Ubicacion,
				TipoEntrega:      string(a.TipoEntrega),
				CostoEnvio:       a.CostoEnvio,
				ZonaEntrega:      a.ZonaEntrega,
				Alias:            a.Alias,
				Cbu:              a.Cbu,
				DireccionXmr:     a.DireccionXmr,
				Cantidad:         a.Cantidad,
			}
		}
		views.Pagina(w, r, views.Compra(data, entregaSlice, total, totalXmr)).Render(r.Context(), w)
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			notificar(w, r, logic.NotificationError, "error", "Formulario inválido", "No pudimos procesar el envío.", "/compra", "Reintentar")
			http.Redirect(w, r, "/compra", http.StatusSeeOther)
			return
		}
		err := h.datos.ComprarCarrito(int32(id))
		switch {
		case errors.Is(err, logic.ErrCarritoVacio):
			notificar(w, r, logic.NotificationInfo, "información", "Tu carrito está vacío", "Agregá una pieza antes de confirmar la compra.", "/", "Seguir explorando")
			http.Redirect(w, r, "/carrito", http.StatusSeeOther)
			return
		case err != nil && strings.Contains(err.Error(), "stock"):
			notificar(w, r, logic.NotificationWarning, "stock insuficiente", "Stock insuficiente", err.Error(), "/carrito", "Ajustar carrito")
			http.Redirect(w, r, "/compra", http.StatusSeeOther)
			return
		case err != nil:
			notificar(w, r, logic.NotificationError, "error", "No pudimos confirmar la compra", "Ocurrió un error inesperado. Intentá de nuevo.", "/compra", "Reintentar")
			http.Error(w, "Error al confirmar la compra: "+err.Error(), http.StatusInternalServerError)
			return
		}
		notificar(w, r, logic.NotificationSuccess, "éxito", "Compra confirmada", "Tu pedido quedó registrado y el stock se actualizó.", "/mi_cuenta", "Ver mis pedidos")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	default:
		return
	}
}

func (h *ArticuloHandler) FavoritosHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	arts, err := h.datos.Queries().GetFavoritosByUsuario(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
		return
	}
	if categoria := r.FormValue("categoria"); categoria != "" {
		filtered := make([]db.GetFavoritosByUsuarioRow, 0, len(arts))
		for _, a := range arts {
			if string(a.Categoria) == categoria {
				filtered = append(filtered, a)
			}
		}
		arts = filtered
	}
	data := make([]views.ArticuloData, len(arts))
	for i, a := range arts {
		data[i] = views.ArticuloData{
			IDArticulo:       a.IDArticulo,
			Titulo:           a.Titulo,
			Descripcion:      a.Descripcion,
			Condicion:        string(a.Condicion),
			Categoria:        string(a.Categoria),
			RutaImagen:       a.RutaImagen,
			Precio:           a.Precio,
			PrecioXmr:        a.PrecioXmr,
			Stock:            a.Stock,
			FechaPublicacion: a.FechaPublicacion,
			Ubicacion:        a.Ubicacion,
			TipoEntrega:      string(a.TipoEntrega),
			CostoEnvio:       a.CostoEnvio,
			ZonaEntrega:      a.ZonaEntrega,
			Alias:            "",
			Cbu:              "",
			DireccionXmr:     sql.NullString{},
			Cantidad:         0,
		}
	}
	if logic.IsHtmxRequest(r) {
		views.Pagina(w, r, views.GrillaFavoritos(data)).Render(r.Context(), w)
		return
	}
	views.Pagina(w, r, views.Favoritos(data, views.BuildFiltrosFav(arts))).Render(r.Context(), w)
}

func (h *ArticuloHandler) AddCarritoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userId, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	artID, err := strconv.Atoi(r.FormValue("id_articulo"))
	if err != nil {
		http.Error(w, "ID de articulo inválido", http.StatusBadRequest)
		return
	}
	if err := h.datos.Queries().AddToCarrito(r.Context(), db.AddToCarritoParams{
		IDUsuario:  int32(userId),
		IDArticulo: int32(artID),
		Cantidad:   1,
	}); err != nil {
		notificar(w, r, logic.NotificationError, "error", "No pudimos agregar al carrito", "Verificá que la pieza todavía tenga stock.", "/", "Seguir explorando")
		http.Redirect(w, r, destinoPost(r), http.StatusSeeOther)
		return
	}
	notificar(w, r, logic.NotificationInfo, "información", "Agregado al carrito", "La pieza quedó reservada en tu carrito.", "/carrito", "Ver carrito")
	http.Redirect(w, r, destinoPost(r), http.StatusSeeOther)
}

func (h *ArticuloHandler) AddFavoritoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userId, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	artID, err := strconv.Atoi(r.FormValue("id_articulo"))
	if err != nil {
		http.Error(w, "ID de articulo inválido", http.StatusBadRequest)
		return
	}
	if err := h.datos.Queries().AddFavorito(r.Context(), db.AddFavoritoParams{
		IDUsuario:  int32(userId),
		IDArticulo: int32(artID),
	}); err != nil {
		notificar(w, r, logic.NotificationError, "error", "No pudimos guardar el favorito", "Intentá de nuevo en unos segundos.", "/favoritos", "Ver favoritos")
		http.Redirect(w, r, destinoPost(r), http.StatusSeeOther)
		return
	}
	notificar(w, r, logic.NotificationSuccess, "éxito", "Agregado a favoritos", "Guardaste la pieza en tu colección.", "/favoritos", "Ver favoritos")
	http.Redirect(w, r, destinoPost(r), http.StatusSeeOther)
}

// destinoPost devuelve el origen del POST para el patrón PRG.
func destinoPost(r *http.Request) string {
	if ref := r.Header.Get("Referer"); ref != "" {
		return ref
	}
	return "/"
}

// DELETE /articulos/{id} - Eliminar articulo publicado
func (h *ArticuloHandler) DeleteArticuloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	userId, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/articulos/"))
	if err != nil {
		http.Error(w, "ID de articulo inválido", http.StatusBadRequest)
		return
	}
	filas, err := h.datos.Queries().DeleteArticuloDeVendedor(r.Context(), db.DeleteArticuloDeVendedorParams{
		IDArticulo: int32(id),
		IDVendedor: int32(userId),
	})
	if err != nil {
		http.Error(w, "Error al eliminar el artículo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		http.Error(w, "No podés eliminar este artículo", http.StatusNotFound)
		return
	}
	notificar(w, r, logic.NotificationSuccess, "éxito", "Artículo eliminado", "La publicación se borró correctamente.", "/", "Volver al inicio")
	responderNotificaciones(w, r)
}

// responderNotificaciones devuelve el contenedor de notificaciones con swap OOB
// para que htmx lo actualice sin recargar la página.
func responderNotificaciones(w http.ResponseWriter, r *http.Request) {
	session := logic.SessionID(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Notificaciones(logic.TakeNotifications(session)).Render(r.Context(), w)
}

// DELETE /carrito/{id} - Eliminar articulo del carrito
func (h *ArticuloHandler) RemoveCarritoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/carrito/"))
	if err != nil {
		http.Error(w, "ID de articulo inválido", http.StatusBadRequest)
		return
	}
	userId, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	err = h.datos.Queries().RemoveFromCarrito(r.Context(), db.RemoveFromCarritoParams{
		IDUsuario:  int32(userId),
		IDArticulo: int32(id),
	})
	if err != nil {
		http.Error(w, "Error al eliminar del carrito: "+err.Error(), http.StatusInternalServerError)
		return
	}
	notificar(w, r, logic.NotificationInfo, "información", "Pieza quitada del carrito", "La pieza ya no ocupa lugar en tu carrito.", "/carrito", "Ver carrito")
	responderNotificaciones(w, r)
}

// DELETE /favoritos/{id} - Eliminar favorito
func (h *ArticuloHandler) RemoveFavoritoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/favoritos/"))
	if err != nil {
		http.Error(w, "ID de articulo inválido", http.StatusBadRequest)
		return
	}
	userId, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	err = h.datos.Queries().RemoveFavorito(r.Context(), db.RemoveFavoritoParams{
		IDUsuario:  int32(userId),
		IDArticulo: int32(id),
	})
	if err != nil {
		http.Error(w, "Error al eliminar el favorito: "+err.Error(), http.StatusInternalServerError)
		return
	}
	notificar(w, r, logic.NotificationInfo, "información", "Quitaste la pieza de favoritos", "La pieza salió de tu colección.", "/favoritos", "Ver favoritos")
	responderNotificaciones(w, r)
}

// SearchHandler responde el fragmento de resultados para el buscador del header.
func (h *ArticuloHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	arts, err := h.datos.Queries().SearchArticulos(r.Context(), sql.NullString{String: q, Valid: q != ""})
	if err != nil {
		http.Error(w, "Error al buscar articulos", http.StatusInternalServerError)
		return
	}
	data := make([]views.ArticuloData, len(arts))
	for i, a := range arts {
		data[i] = views.ArticuloData{
			IDArticulo:       a.IDArticulo,
			Titulo:           a.Titulo,
			Descripcion:      a.Descripcion,
			Condicion:        string(a.Condicion),
			Categoria:        string(a.Categoria),
			RutaImagen:       a.RutaImagen,
			Precio:           a.Precio,
			PrecioXmr:        a.PrecioXmr,
			Stock:            a.Stock,
			FechaPublicacion: a.FechaPublicacion,
			Ubicacion:        a.Ubicacion,
			TipoEntrega:      string(a.TipoEntrega),
			CostoEnvio:       a.CostoEnvio,
			ZonaEntrega:      a.ZonaEntrega,
		}
	}
	views.Pagina(w, r, views.Busqueda(data)).Render(r.Context(), w)
}

// FichaArticuloHandler devuelve el contenido del modal de ficha (htmx GET).
func (h *ArticuloHandler) FichaArticuloHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/articulos/"))
	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	row, err := h.datos.Queries().GetArticuloCompletoByID(r.Context(), int32(id))
	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	data := views.ArticuloData{
		IDArticulo:       row.IDArticulo,
		Titulo:           row.Titulo,
		Descripcion:      row.Descripcion,
		Condicion:        string(row.Condicion),
		Categoria:        string(row.Categoria),
		RutaImagen:       row.RutaImagen,
		Precio:           row.Precio,
		PrecioXmr:        row.PrecioXmr,
		Stock:            row.Stock,
		FechaPublicacion: row.FechaPublicacion,
		Ubicacion:        row.Ubicacion,
		TipoEntrega:      string(row.TipoEntrega),
		CostoEnvio:       row.CostoEnvio,
		ZonaEntrega:      row.ZonaEntrega,
		Alias:            row.Alias,
		Cbu:              row.Cbu,
		DireccionXmr:     row.DireccionXmr,
	}
	views.Pagina(w, r, views.VistaFicha(data)).Render(r.Context(), w)
}

// EditarArticuloHandler devuelve el modal con el formulario de edición.
func (h *ArticuloHandler) EditarArticuloHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	id, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/articulos/"), "/editar"))
	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	row, err := h.datos.Queries().GetArticuloCompletoByVendedor(r.Context(), db.GetArticuloCompletoByVendedorParams{
		IDArticulo: int32(id),
		IDVendedor: int32(userId),
	})
	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	data := views.ArticuloData{
		IDArticulo:       row.IDArticulo,
		Titulo:           row.Titulo,
		Descripcion:      row.Descripcion,
		Condicion:        string(row.Condicion),
		Categoria:        string(row.Categoria),
		RutaImagen:       row.RutaImagen,
		Precio:           row.Precio,
		PrecioXmr:        row.PrecioXmr,
		Stock:            row.Stock,
		FechaPublicacion: row.FechaPublicacion,
		Ubicacion:        row.Ubicacion,
		TipoEntrega:      string(row.TipoEntrega),
		CostoEnvio:       row.CostoEnvio,
		ZonaEntrega:      row.ZonaEntrega,
		Alias:            row.Alias,
		Cbu:              row.Cbu,
		DireccionXmr:     row.DireccionXmr,
	}
	views.Pagina(w, r, views.VistaEditarArticulo(data)).Render(r.Context(), w)
}

// ActualizarArticuloHandler aplica la edición vía sqlc.
func (h *ArticuloHandler) ActualizarArticuloHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/articulos/"))
	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		notificar(w, r, logic.NotificationError, "error", "Formulario inválido", "No pudimos procesar el envío.", "/mi_cuenta", "Reintentar")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	in, ok := logic.ParseActualizarArticulo(r)
	if !ok {
		notificar(w, r, logic.NotificationError, "error", "No pudimos guardar los cambios", "Revisá título, descripción, categoría, condición, precio, stock, ubicación y entrega.", "/mi_cuenta", "Reintentar")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	filas, err := h.datos.ActualizarArticulo(int32(id), int32(userId), in)
	if err != nil {
		notificar(w, r, logic.NotificationError, "error", "No pudimos guardar los cambios", "Ocurrió un error al actualizar la publicación.", "/mi_cuenta", "Reintentar")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	if filas == 0 {
		notificar(w, r, logic.NotificationWarning, "atención", "No podés editar esta pieza", "La publicación no existe o no te pertenece.", "/", "Volver al inicio")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	notificar(w, r, logic.NotificationSuccess, "éxito", "Artículo actualizado", in.Titulo, "/", "Ver publicación")
	http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
}
