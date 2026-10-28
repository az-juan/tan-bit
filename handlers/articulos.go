package handlers

import (
	"encoding/json"
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
	if err != nil {
		http.Error(w, "Articulo no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articulo)
}

// PUT /articulos/{id} - Actualizar articulo
func (h *ArticuloHandler) updateArticulo(w http.ResponseWriter, r *http.Request, id int32) {
	/*var updatedArticulo db.Articulo

	err := json.NewDecoder(r.Body).Decode(&updatedArticulo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	index, err := findArticuloIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Mantener el ID original
	updatedArticulo.IDArticulo = int32(id)
	articulos[index] = updatedArticulo
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedArticulo)*/
}

// DELETE /articulos/{id} - Eliminar articulo
func (h *ArticuloHandler) deleteArticulo(w http.ResponseWriter, r *http.Request, id int32) {
	/*index, err := findArticuloIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	articulos = append(articulos[:index], articulos[index+1:]...)
	w.WriteHeader(http.StatusNoContent)*/
}

/*
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sortColumn := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")
	if sortOrder != "desc" {
		sortOrder = "asc" // Valor por defecto
	}
	var users []db.Articulo
	var err error
	switch sortColumn {
	case "id":
		if sortOrder == "asc" {
			users, err = h.queries.ListUsersOrderByIdAsc(ctx)
		} else {
			users, err = h.queries.ListUsersOrderByIdDesc(ctx)
		}
	case "name":
		if sortOrder == "asc" {
			users, err = h.queries.ListUsersOrderByNameAsc(ctx)
		} else {
			users, err = h.queries.ListUsersOrderByNameDesc(ctx)
		}
	case "email":
		if sortOrder == "asc" {
			users, err = h.queries.ListUsersOrderByEmailAsc(ctx)
		} else {
			users, err = h.queries.ListUsersOrderByEmailDesc(ctx)
		}
	default:
		// Orden por defecto
		users, err = h.queries.ListUsers(ctx)
		sortColumn = "id" // Columna por defecto para resaltar
		sortOrder = "asc"
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	views.UserPage(users, sortColumn, sortOrder).Render(r.Context(), w)
}
*/
