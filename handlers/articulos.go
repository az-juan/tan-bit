package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "tan-bit.com/tan-bit/db/sqlc"
)

type ArticuloHandler struct {
	queries *db.Queries
}

type CreateArticuloRequest struct {
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Condicion   string `json:"condicion"`
	Categoria   string `json:"categoria"`
	Precio      string `json:"precio"`
	Stock       int32  `json:"stock"`
	IDVendedor  int32  `json:"id_vendedor"`
	Ubicacion   string `json:"ubicacion"`
}

type UpdateArticuloRequest struct {
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Condicion   string `json:"condicion"`
	Categoria   string `json:"categoria"`
	Precio      string `json:"precio"`
	Stock       int32  `json:"stock"`
}

type ArticuloResponse struct {
	IDArticulo       int32      `json:"id_articulo"`
	Titulo           string     `json:"titulo"`
	Descripcion      string     `json:"descripcion"`
	Condicion        string     `json:"condicion"`
	Categoria        string     `json:"categoria"`
	Precio           string     `json:"precio"`
	Stock            int32      `json:"stock"`
	IDVendedor       int32      `json:"id_vendedor,omitempty"`
	FechaPublicacion *time.Time `json:"fecha_publicacion,omitempty"`
}

func toArticuloFromCreate(a db.CreateArticuloRow) ArticuloResponse {
	return ArticuloResponse{
		IDArticulo:  a.IDArticulo,
		Titulo:      a.Titulo,
		Descripcion: a.Descripcion,
		Condicion:   string(a.Condicion),
		Categoria:   string(a.Categoria),
		Precio:      a.Precio,
		Stock:       a.Stock,
		IDVendedor:  a.IDVendedor,
	}
}

func toArticuloFromGetByID(a db.GetArticuloByIDRow) ArticuloResponse {
	var fecha *time.Time
	if a.FechaPublicacion.Valid {
		fecha = &a.FechaPublicacion.Time
	}
	return ArticuloResponse{
		IDArticulo:       a.IDArticulo,
		Titulo:           a.Titulo,
		Descripcion:      a.Descripcion,
		Condicion:        string(a.Condicion),
		Categoria:        string(a.Categoria),
		Precio:           a.Precio,
		Stock:            a.Stock,
		FechaPublicacion: fecha,
	}
}

func toArticulosFromGet(arts []db.GetArticulosRow) []ArticuloResponse {
	res := make([]ArticuloResponse, 0, len(arts))
	for _, a := range arts {
		var fecha *time.Time
		if a.FechaPublicacion.Valid {
			fecha = &a.FechaPublicacion.Time
		}
		res = append(res, ArticuloResponse{
			IDArticulo:       a.IDArticulo,
			Titulo:           a.Titulo,
			Descripcion:      a.Descripcion,
			Condicion:        string(a.Condicion),
			Categoria:        string(a.Categoria),
			Precio:           a.Precio,
			Stock:            a.Stock,
			FechaPublicacion: fecha,
		})
	}
	return res
}

func NewArticuloHandler(queries *db.Queries) *ArticuloHandler {
	return &ArticuloHandler{
		queries: queries,
	}
}

// Manejador para /articulos
func (h *ArticuloHandler) ArtsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// case http.MethodGet:
	// h.GetArticulos(w, r)
	case http.MethodPost:
		h.createArticulo(w, r)
	default:
		http.Error(w, "Metodo no permitido",
			http.StatusMethodNotAllowed)
	}
}

// Manejador para /articulos/{id}
func (h *ArticuloHandler) ArtHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "ID de articulo invalido", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getArticulo(w, r, int32(id))
	case http.MethodPut:
		h.updateArticulo(w, r, int32(id))
	case http.MethodDelete:
		h.deleteArticulo(w, r, int32(id))
	default:
		http.Error(w, "Metodo no permitido",
			http.StatusMethodNotAllowed)
	}
	// templ.Handler(views.UserPage(models.Users)).ServeHTTP(w, r)
}

// GET /articulos - Obtener todos los articulos
// func (h *ArticuloHandler) GetArticulos(w http.ResponseWriter, r *http.Request) {
// 	arts, err := h.queries.GetArticulos(r.Context())
// 	if err != nil {
// 		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
// 		return
// 	}
// 	views.Marketplace(arts).Render(r.Context(), w)
// }

// POST /articulos - Crear nuevo articulo
func (h *ArticuloHandler) createArticulo(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" && (r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" || (len(r.Header.Get("Content-Type")) >= 33 && r.Header.Get("Content-Type")[:33] == "multipart/form-data")) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulario inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
		titulo := r.FormValue("titulo")
		descripcion := r.FormValue("descripcion")
		condicion := r.FormValue("condicion")
		categoria := r.FormValue("categoria")
		precio := r.FormValue("precio")
		stockStr := r.FormValue("stock")
		vendedorStr := r.FormValue("id_vendedor")
		stock := int32(1)
		if v, err := strconv.Atoi(stockStr); err == nil {
			stock = int32(v)
		}
		idVendedor := int32(1)
		if v, err := strconv.Atoi(vendedorStr); err == nil {
			idVendedor = int32(v)
		}
		params := db.CreateArticuloParams{
			Titulo:      titulo,
			Descripcion: descripcion,
			Condicion:   db.Condicion(condicion),
			Categoria:   db.Categoria(categoria),
			Precio:      precio,
			Stock:       stock,
			IDVendedor:  idVendedor,
			Ubicacion:   r.FormValue("ubicacion"),
		}
		_, err := h.queries.CreateArticulo(r.Context(), params)
		if err != nil {
			http.Error(w, "Error al crear el artículo: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var req CreateArticuloRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	params := db.CreateArticuloParams{
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Condicion:   db.Condicion(req.Condicion),
		Categoria:   db.Categoria(req.Categoria),
		Precio:      req.Precio,
		Stock:       req.Stock,
		IDVendedor:  req.IDVendedor,
		Ubicacion:   req.Ubicacion,
	}
	nuevoArticulo, err := h.queries.CreateArticulo(r.Context(), params)
	if esViolacionFK(err) {
		http.Error(w, "No existe tal vendedor", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Error al crear el artículo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toArticuloFromCreate(nuevoArticulo))
}

// GET /articulos/{id} - Obtener articulo específico
func (h *ArticuloHandler) getArticulo(w http.ResponseWriter, r *http.Request, id int32) {
	articulo, err := h.queries.GetArticuloByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Articulo no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener el articulo", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toArticuloFromGetByID(articulo))
}

// PUT /articulos/{id} - Actualizar articulo
func (h *ArticuloHandler) updateArticulo(w http.ResponseWriter, r *http.Request, id int32) {
	var req UpdateArticuloRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	params := db.UpdateArticuloParams{
		IDArticulo:  id,
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Condicion:   db.Condicion(req.Condicion),
		Categoria:   db.Categoria(req.Categoria),
		Precio:      req.Precio,
		Stock:       req.Stock,
	}

	filas, err := h.queries.UpdateArticulo(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar el artículo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		http.Error(w, "Articulo no encontrado", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /articulos/{id} - Eliminar articulo
func (h *ArticuloHandler) deleteArticulo(w http.ResponseWriter, r *http.Request, id int32) {
	filas, err := h.queries.DeleteArticulo(r.Context(), id)
	if esViolacionFK(err) {
		http.Error(w, "El articulo tiene compras asociadas", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Error al eliminar el artículo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		http.Error(w, "Articulo no encontrado", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
