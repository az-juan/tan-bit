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

type VendedorHandler struct {
	queries *db.Queries
}

func NewVendedorHandler(queries *db.Queries) *VendedorHandler {
	return &VendedorHandler{
		queries: queries,
	}
}

// Manejador para /vendedores
func (h *VendedorHandler) VendsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getVendedores(w, r)
	case http.MethodPost:
		h.createVendedor(w, r)
	default:
		http.Error(w, "Metodo no permitido",
			http.StatusMethodNotAllowed)
	}
}

// Manejador para /vendedores/{id}
func (h *VendedorHandler) VendHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "ID de vendedor invalido", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getVendedor(w, r, int32(id))
	case http.MethodPut:
		h.updateVendedor(w, r, int32(id))
	case http.MethodDelete:
		h.deleteVendedor(w, r, int32(id))
	default:
		http.Error(w, "Metodo no permitido",
			http.StatusMethodNotAllowed)
	}
}

// GET /vendedores - Obtener todos los vendedores
func (h *VendedorHandler) getVendedores(w http.ResponseWriter, r *http.Request) {
	vends, err := h.queries.GetVendedores(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener vendedores", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(vends)
}

// POST /vendedores - Crear nuevo vendedor
func (h *VendedorHandler) createVendedor(w http.ResponseWriter, r *http.Request) {
	var params db.CreateVendedorParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "JSON inválido: "+ err.Error(), http.StatusBadRequest)
		return
	}

	nuevoVendedor, err := h.queries.CreateVendedor(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear el vendedor: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoVendedor)
}


// GET /vendedores/{id} - Obtener vendedor específico
func (h *VendedorHandler) getVendedor(w http.ResponseWriter, r *http.Request, id int32) {
	vendedor, err := h.queries.GetVendedorByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Vendedor no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener el vendedor", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendedor)
}

// PUT /vendedores/{id} - Actualizar vendedor
func (h *VendedorHandler) updateVendedor(w http.ResponseWriter, r *http.Request, id int32) {
	var params db.UpdateVendedorParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "JSON inválido: "+ err.Error(), http.StatusBadRequest)
		return
	}
	params.IDVendedor = id

	filas, err := h.queries.UpdateVendedor(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar el vendedor: "+ err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		http.Error(w, "Vendedor no encontrado", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /vendedores/{id} - Eliminar vendedor
func (h *VendedorHandler) deleteVendedor(w http.ResponseWriter, r *http.Request, id int32) {
	filas, err := h.queries.DeleteVendedor(r.Context(), id)
	if esViolacionFK(err) {
		http.Error(w, "El vendedor tiene articulos publicados", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Error al eliminar el vendedor: "+ err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		http.Error(w, "Vendedor no encontrado", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
