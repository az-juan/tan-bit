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

type UsuarioHandler struct {
	queries *db.Queries
}

func NewUsuarioHandler(queries *db.Queries) *UsuarioHandler {
	return &UsuarioHandler{
		queries: queries,
	}
}

// Manejador para /usuarios
func (h *UsuarioHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUsuarios(w, r)
	case http.MethodPost:
		h.createUsuario(w, r)
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

// Manejador para /usuarios/{id}
func (h *UsuarioHandler) UserHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "ID de usuario invalido", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getUsuario(w, r, int32(id))
	case http.MethodPut:
		h.updateUsuario(w, r, int32(id))
	case http.MethodDelete:
		h.deleteUsuario(w, r, int32(id))
	default:
		http.Error(w, "Metodo no permitido",
			http.StatusMethodNotAllowed)
	}
}

// GET /usuarios - Obtener todos los usuarios
func (h *UsuarioHandler) getUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.queries.GetUsuarios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

// POST /usuarios - Crear nuevo usuario
func (h *UsuarioHandler) createUsuario(w http.ResponseWriter, r *http.Request) {
	var params db.CreateUsuarioParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	nuevoUsuario, err := h.queries.CreateUsuario(r.Context(), params)
	if esViolacionUnique(err) {
		http.Error(w, "Email ya existente.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Error al crear el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoUsuario)
}

// GET /usuarios/{id} - Obtener usuario específico
func (h *UsuarioHandler) getUsuario(w http.ResponseWriter, r *http.Request, id int32) {
	usuario, err := h.queries.GetUsuarioByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Usuario no encontrado.", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// PUT /usuarios/{id} - Actualizar usuario
func (h *UsuarioHandler) updateUsuario(w http.ResponseWriter, r *http.Request, id int32) {
	var params db.UpdateUsuarioParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	params.IDUsuario = id

	filas, err := h.queries.UpdateUsuario(r.Context(), params)
	if esViolacionUnique(err) {
		http.Error(w, "Email ya existente.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Error al actualizar el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /usuarios/{id} - Eliminar usuario
func (h *UsuarioHandler) deleteUsuario(w http.ResponseWriter, r *http.Request, id int32) {
	filas, err := h.queries.DeleteUsuario(r.Context(), id)
	if esViolacionFK(err) {
		http.Error(w, "El usuario es vendedor o tiene compras asociadas, no es posible eliminarlo.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Error al eliminar el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if filas == 0 {
		http.Error(w, "Usuario no encontrado.", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
