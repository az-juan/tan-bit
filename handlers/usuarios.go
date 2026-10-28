package handlers
    
import (
	"encoding/json"
	"net/http"
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
	if err != nil {
		http.Error(w, "Error al crear el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoUsuario)
}