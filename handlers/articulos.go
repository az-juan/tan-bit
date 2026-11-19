package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	db "tan-bit.com/tan-bit/db/sqlc"
	"tan-bit.com/tan-bit/logic"
	"tan-bit.com/tan-bit/views"
)

type ArticuloHandler struct {
	queries *db.Queries
}

func NewArticuloHandler(queries *db.Queries) *ArticuloHandler {
	return &ArticuloHandler{
		queries: queries,
	}
}

func (h *ArticuloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.MarketplaceHandler(w, r)
	case "/carrito":
		h.CarritoHandler(w, r)
	case "/carrito/":
		h.RemoveCarritoHandler(w, r)
	case "/compra":
		h.CompraHandler(w, r)
	case "/favoritos":
		h.FavoritosHandler(w, r)
	case "/favoritos/":
		h.RemoveFavoritoHandler(w, r)
	case "/articulos/":
		h.DeleteArticuloHandler(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (h *ArticuloHandler) MarketplaceHandler(w http.ResponseWriter, r *http.Request) {
	categoria := r.FormValue("categoria")
	var arts []db.GetArticulosRow
	var err error
	if categoria != "" {
		artsRows, err := h.queries.GetArticulosByCategoria(r.Context(), db.Categoria(categoria))
		if err != nil {
			http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
			return
		}
		arts = make([]db.GetArticulosRow, len(artsRows))
		for i, a := range artsRows {
			arts[i] = db.GetArticulosRow{
				IDArticulo:       a.IDArticulo,
				Titulo:           a.Titulo,
				Descripcion:      a.Descripcion,
				Condicion:        a.Condicion,
				Categoria:        a.Categoria,
				Precio:           a.Precio,
				PrecioXmr:        a.PrecioXmr,
				Stock:            a.Stock,
				FechaPublicacion: a.FechaPublicacion,
				RutaImagen:       a.RutaImagen,
				Ubicacion:        a.Ubicacion,
				TipoEntrega:      a.TipoEntrega,
				CostoEnvio:       a.CostoEnvio,
				ZonaEntrega:      a.ZonaEntrega,
			}
		}
	} else {
		arts, err = h.queries.GetArticulos(r.Context())
	}
	if err != nil {
		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
		return
	}
	views.Pagina(r, views.Marketplace(arts, views.BuildFiltros(arts))).Render(r.Context(), w)
}

func (h *ArticuloHandler) CarritoHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	arts, err := h.queries.GetCarritoByUsuario(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
		return
	}
	var total float64
	var totalXmr float64
	for _, a := range arts {
		total += a.Precio * float64(a.Cantidad)
		totalXmr += a.PrecioXmr * float64(a.Cantidad)
	}
	views.Pagina(r, views.Carrito(arts, total, totalXmr)).Render(r.Context(), w)
}

func (h *ArticuloHandler) CompraHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	if r.Method == http.MethodGet {
		arts, err := h.queries.GetCarritoCompra(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
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
		for _, a := range arts {
			total += a.Precio * float64(a.Cantidad)
			totalXmr += a.PrecioXmr * float64(a.Cantidad)
		}
		views.Pagina(r, views.Compra(arts, entregaSlice, total, totalXmr)).Render(r.Context(), w)
		return
	}
	// POST
	arts, err := h.queries.GetCarritoCompra(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
		return
	}
	var total float64
	var totalXmr float64
	for _, a := range arts {
		total += a.Precio * float64(a.Cantidad)
		totalXmr += a.PrecioXmr * float64(a.Cantidad)
	}
	_ = r.FormValue("entrega")
	compra, err := h.queries.CreateCompra(r.Context(), db.CreateCompraParams{
		IDComprador: int32(id),
		Total:       total,
		TotalXmr:    totalXmr,
	})
	if err != nil {
		http.Error(w, "Error al crear compra", http.StatusInternalServerError)
		return
	}
	for _, a := range arts {
		h.queries.CreateDetalleCompra(r.Context(), db.CreateDetalleCompraParams{
			IDCompra:                   compra.IDCompra,
			IDArticulo:                 a.IDArticulo,
			Cantidad:                   a.Cantidad,
			PrecioHistoricoUnitario:    a.Precio,
			PrecioHistoricoUnitarioXmr: a.PrecioXmr,
		})
	}
	h.queries.ClearCarrito(r.Context(), int32(id))
	views.Pagina(r, views.MiCuenta("", "", "", "", "", "", "", []db.GetArticulosByVendedorRow{}, []db.GetFavoritosByUsuarioRow{}, []db.GetComprasByUsuarioRow{})).Render(r.Context(), w)
}

func (h *ArticuloHandler) FavoritosHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	arts, err := h.queries.GetFavoritosByUsuario(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
		return
	}
	views.Pagina(r, views.Favoritos(arts, views.BuildFiltrosFav(arts))).Render(r.Context(), w)
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
	artID, _ := strconv.Atoi(r.FormValue("id_articulo"))
	h.queries.AddToCarrito(r.Context(), db.AddToCarritoParams{
		IDUsuario:  int32(userId),
		IDArticulo: int32(artID),
		Cantidad:   1,
	})
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
	artID, _ := strconv.Atoi(r.FormValue("id_articulo"))
	h.queries.AddFavorito(r.Context(), db.AddFavoritoParams{
		IDUsuario:  int32(userId),
		IDArticulo: int32(artID),
	})
}

// DELETE /articulos/{id} - Eliminar articulo publicado
func (h *ArticuloHandler) DeleteArticuloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/articulos/"))
	if err != nil {
		http.Error(w, "ID de articulo inválido", http.StatusBadRequest)
		return
	}
	filas, err := h.queries.DeleteArticulo(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar el artículo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Artículo no encontrado")
		return
	}
	w.WriteHeader(http.StatusOK)
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
	err = h.queries.RemoveFromCarrito(r.Context(), db.RemoveFromCarritoParams{
		IDUsuario:  int32(userId),
		IDArticulo: int32(id),
	})
	if err != nil {
		http.Error(w, "Error al eliminar del carrito: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	err = h.queries.RemoveFavorito(r.Context(), db.RemoveFavoritoParams{
		IDUsuario:  int32(userId),
		IDArticulo: int32(id),
	})
	if err != nil {
		http.Error(w, "Error al eliminar el favorito: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
