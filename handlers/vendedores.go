package handlers

import (
	"encoding/json"
	"net/http"
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