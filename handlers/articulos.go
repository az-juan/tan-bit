package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	db "tan-bit.com/tan-bit/db/sqlc"
)

type ArticuloHandler struct {
	queries *db.Queries
}

func NewArticuloHandler(queries *db.Queries) *ArticuloHandler {
	return &ArticuloHandler{
		queries: queries,
	}
}

// Manejador para /articulos
func (h *ArticuloHandler) ArtsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getArticulos(w, r)
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
}

// GET /articulos - Obtener todos los articulo
func (h *ArticuloHandler) getArticulos(w http.ResponseWriter, r *http.Request) {
	arts, err := h.queries.GetArticulos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener articulos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(arts)
}

// POST /articulos - Crear nuevo articulo
func (h *ArticuloHandler) createArticulo(w http.ResponseWriter, r *http.Request) {
	var params db.CreateArticuloParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "JSON inválido: "+ err.Error(), http.StatusBadRequest)
		return
	}

	nuevoArticulo, err := h.queries.CreateArticulo(r.Context(), params)
	if esViolacionFK(err) {
		http.Error(w, "No existe tal vendedor", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Error al crear el artículo: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoArticulo)
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
	json.NewEncoder(w).Encode(articulo)
}

// PUT /articulos/{id} - Actualizar articulo
func (h *ArticuloHandler) updateArticulo(w http.ResponseWriter, r *http.Request, id int32) {
	var params db.UpdateArticuloParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "JSON inválido: "+ err.Error(), http.StatusBadRequest)
		return
	}
	params.IDArticulo = id

	filas, err := h.queries.UpdateArticulo(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar el artículo: "+ err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Error al eliminar el artículo: "+ err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		http.Error(w, "Articulo no encontrado", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
