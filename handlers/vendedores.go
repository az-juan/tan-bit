package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	db "tan-bit.com/tan-bit/db/sqlc"
	"tan-bit.com/tan-bit/views"
)

type VendedorHandler struct {
	queries *db.Queries
}

type CreateVendedorRequest struct {
	IDVendedor int32 `json:"id_vendedor"`
	Telefono   int64 `json:"telefono"`
}

type UpdateVendedorRequest struct {
	Telefono int64 `json:"telefono"`
}

type VendedorResponse struct {
	IDVendedor int32 `json:"id_vendedor"`
	Telefono   int64 `json:"telefono"`
}

func toVendedorResponse(v db.Vendedor) VendedorResponse {
	return VendedorResponse{
		IDVendedor: v.IDVendedor,
		Telefono:   v.Telefono,
	}
}

func toVendedoresFromGet(vends []db.Vendedor) []VendedorResponse {
	res := make([]VendedorResponse, 0, len(vends))
	for _, v := range vends {
		res = append(res, toVendedorResponse(v))
	}
	return res
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
	json.NewEncoder(w).Encode(toVendedoresFromGet(vends))
}

// POST /vendedores - Crear nuevo vendedor
func (h *VendedorHandler) createVendedor(w http.ResponseWriter, r *http.Request) {
	var req CreateVendedorRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	params := db.CreateVendedorParams{
		IDVendedor: req.IDVendedor,
		Telefono:   req.Telefono,
	}

	nuevoVendedor, err := h.queries.CreateVendedor(r.Context(), params)
	if esViolacionFK(err) {
		http.Error(w, "No existe tal usuario", http.StatusBadRequest)
		return
	}
	if esViolacionUnique(err) {
		http.Error(w, "El usuario ya es vendedor", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Error al crear el vendedor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toVendedorResponse(nuevoVendedor))
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
	json.NewEncoder(w).Encode(toVendedorResponse(vendedor))
}

// PUT /vendedores/{id} - Actualizar vendedor
func (h *VendedorHandler) updateVendedor(w http.ResponseWriter, r *http.Request, id int32) {
	var req UpdateVendedorRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	params := db.UpdateVendedorParams{
		IDVendedor: id,
		Telefono:   req.Telefono,
	}

	filas, err := h.queries.UpdateVendedor(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar el vendedor: "+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Error al eliminar el vendedor: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		http.Error(w, "Vendedor no encontrado", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Manejador de página de venta
func (h *VendedorHandler) VenderHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id, _ := strconv.Atoi(cookie.Value)
	v, err := h.queries.GetVendedorByID(r.Context(), int32(id))
	if err != nil || v.Telefono == 0 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<script>alert('Debe cargar un telefono para vender'); window.location='/mi_cuenta';</script>"))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Vender(views.VenderData{}, views.CategoriasDB(), views.CondicionesDB()).Render(r.Context(), w)
}

// ServeHTTP para vendedor handler
func (h *VendedorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/vender" {
		h.VenderHandler(w, r)
		return
	}
	if r.URL.Path == "/vendedores" || r.URL.Path == "/vendedores/" {
		h.VendsHandler(w, r)
		return
	}
	h.VendHandler(w, r)
}
