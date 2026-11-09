package handlers

import (
	// "database/sql"
	// "encoding/json"
	// "errors"
	"net/http"
	"strconv"
	// "strings"

	db "tan-bit.com/tan-bit/db/sqlc"
	"tan-bit.com/tan-bit/views"
)

type UsuarioHandler struct {
	queries *db.Queries
}

func NewUsuarioHandler(queries *db.Queries) *UsuarioHandler {
	return &UsuarioHandler{
		queries: queries,
	}
}

// type UsuarioRequest struct {
// 	Apellido string `json:"apellido"`
// 	Nombre   string `json:"nombre"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// type UsuarioResponse struct {
// 	IDUsuario int32  `json:"id_usuario"`
// 	Apellido  string `json:"apellido"`
// 	Nombre    string `json:"nombre"`
// 	Email     string `json:"email"`
// }

// func toUsuarioResponse(u db.Usuario) UsuarioResponse {
// 	return UsuarioResponse{
// 		IDUsuario: u.IDUsuario,
// 		Apellido:  u.Apellido,
// 		Nombre:    u.Nombre,
// 		Email:     u.Email,
// 	}
// }

// func toUsuariosFromGet(usuarios []db.Usuario) []UsuarioResponse {
// 	res := make([]UsuarioResponse, 0, len(usuarios))
// 	for _, u := range usuarios {
// 		res = append(res, toUsuarioResponse(u))
// 	}
// 	return res
// }

func (h *UsuarioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		h.LoginHandler(w, r)
		return
	case "/registrarse":
		h.RegistrarseHandler(w, r)
		return
	case "/mi_cuenta":
		h.MiCuentaHandler(w, r)
		return
	default:
		NotFoundHandler(w, r)
	}
}

func (h *UsuarioHandler) RegistrarseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.Registrarse().Render(r.Context(), w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}
	params := db.CreateUsuarioParams{
		Apellido: r.FormValue("apellido"),
		Nombre:   r.FormValue("nombre"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	nuevoUsuario, err := h.queries.CreateUsuario(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if tel := r.FormValue("telefono"); tel != "" {
		t, _ := strconv.ParseInt(tel, 10, 64)
		_, _ = h.queries.CreateVendedor(r.Context(), db.CreateVendedorParams{IDVendedor: nuevoUsuario.IDUsuario, Telefono: t})
	}
	http.SetCookie(w, &http.Cookie{Name: "user_id", Value: strconv.Itoa(int(nuevoUsuario.IDUsuario)), Path: "/"})
	http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
}

// Manejador para /login
func (h *UsuarioHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.Login().Render(r.Context(), w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	u, err := h.queries.GetUsuarioByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}
	if u.Password != password {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "user_id", Value: strconv.Itoa(int(u.IDUsuario)), Path: "/"})
	http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
}

func (h *UsuarioHandler) MiCuentaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	nombre, apellido, email := "", "", ""
	telefono := ""
	if cookie, err := r.Cookie("user_id"); err == nil {
		if id, err := strconv.Atoi(cookie.Value); err == nil {
			if u, err := h.queries.GetUsuarioByID(r.Context(), int32(id)); err == nil {
				nombre = u.Nombre
				apellido = u.Apellido
				email = u.Email
			}
			if v, err := h.queries.GetVendedorByID(r.Context(), int32(id)); err == nil {
				telefono = strconv.FormatInt(v.Telefono, 10)
			}
		}
	}
	views.MiCuenta(nombre, apellido, email, telefono).Render(r.Context(), w)
}

// Manejador para /usuarios
// func (h *UsuarioHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		h.getUsuarios(w, r)
// 	case http.MethodPost:
// 		h.createUsuario(w, r)
// 	default:
// 		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
// 	}
// }

// Manejador para /usuarios/{id}
// func (h *UsuarioHandler) UserHandler(w http.ResponseWriter, r *http.Request) {
// 	// Extraer ID del path
// 	parts := strings.Split(r.URL.Path, "/")
// 	if len(parts) != 3 {
// 		http.Error(w, "Invalid URL", http.StatusBadRequest)
// 		return
// 	}
// 	id, err := strconv.Atoi(parts[2])

// 	if err != nil {
// 		http.Error(w, "ID de usuario invalido", http.StatusBadRequest)
// 		return
// 	}
// 	switch r.Method {
// 	case http.MethodGet:
// 		h.getUsuario(w, r, int32(id))
// 	case http.MethodPut:
// 		h.updateUsuario(w, r, int32(id))
// 	case http.MethodDelete:
// 		h.deleteUsuario(w, r, int32(id))
// 	default:
// 		http.Error(w, "Metodo no permitido",
// 			http.StatusMethodNotAllowed)
// 	}
// }

// GET /usuarios - Obtener todos los usuarios
// func (h *UsuarioHandler) getUsuarios(w http.ResponseWriter, r *http.Request) {
// 	usuarios, err := h.queries.GetUsuarios(r.Context())
// 	if err != nil {
// 		http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(toUsuariosFromGet(usuarios))
// }

// POST /usuarios - Crear nuevo usuario
// func (h *UsuarioHandler) createUsuario(w http.ResponseWriter, r *http.Request) {
// 	ct := r.Header.Get("Content-Type")
// 	if ct == "application/x-www-form-urlencoded" || (len(ct) >= 33 && ct[:33] == "multipart/form-data") {
// 		if err := r.ParseForm(); err != nil {
// 			http.Error(w, "Formulario inválido", http.StatusBadRequest)
// 			return
// 		}
// 		params := db.CreateUsuarioParams{
// 			Apellido: r.FormValue("apellido"),
// 			Nombre:   r.FormValue("nombre"),
// 			Email:    r.FormValue("email"),
// 			Password: r.FormValue("password"),
// 		}
// 		nuevoUsuario, err := h.queries.CreateUsuario(r.Context(), params)
// 		if err != nil {
// 			http.Error(w, "Error al crear el usuario: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		if tel := r.FormValue("telefono"); tel != "" {
// 			t, _ := strconv.ParseInt(tel, 10, 64)
// 			_, _ = h.queries.CreateVendedor(r.Context(), db.CreateVendedorParams{IDVendedor: nuevoUsuario.IDUsuario, Telefono: t})
// 		}
// 		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
// 		return
// 	}

// 	var req UsuarioRequest
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	params := db.CreateUsuarioParams{
// 		Apellido: req.Apellido,
// 		Nombre:   req.Nombre,
// 		Email:    req.Email,
// 		Password: req.Password,
// 	}
// 	nuevoUsuario, err := h.queries.CreateUsuario(r.Context(), params)
// 	if esViolacionUnique(err) {
// 		http.Error(w, "Email ya existente.", http.StatusConflict)
// 		return
// 	}
// 	if err != nil {
// 		http.Error(w, "Error al crear el usuario: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(toUsuarioResponse(nuevoUsuario))
// }

// GET /usuarios/{id} - Obtener usuario específico
// func (h *UsuarioHandler) getUsuario(w http.ResponseWriter, r *http.Request, id int32) {
// 	usuario, err := h.queries.GetUsuarioByID(r.Context(), id)
// 	if errors.Is(err, sql.ErrNoRows) {
// 		http.Error(w, "Usuario no encontrado.", http.StatusNotFound)
// 		return
// 	}
// 	if err != nil {
// 		http.Error(w, "Error al obtener el usuario: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(toUsuarioResponse(usuario))
// }

// PUT /usuarios/{id} - Actualizar usuario
// func (h *UsuarioHandler) updateUsuario(w http.ResponseWriter, r *http.Request, id int32) {
// 	var req UsuarioRequest

// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	params := db.UpdateUsuarioParams{
// 		IDUsuario: id,
// 		Apellido:  req.Apellido,
// 		Nombre:    req.Nombre,
// 		Email:     req.Email,
// 	}

// 	filas, err := h.queries.UpdateUsuario(r.Context(), params)
// 	if esViolacionUnique(err) {
// 		http.Error(w, "Email ya existente.", http.StatusConflict)
// 		return
// 	}
// 	if err != nil {
// 		http.Error(w, "Error al actualizar el usuario: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if filas == 0 {
// 		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
// 		return
// 	}
// 	w.WriteHeader(http.StatusNoContent)
// }

// DELETE /usuarios/{id} - Eliminar usuario
// func (h *UsuarioHandler) deleteUsuario(w http.ResponseWriter, r *http.Request, id int32) {
// 	filas, err := h.queries.DeleteUsuario(r.Context(), id)
// 	if esViolacionFK(err) {
// 		http.Error(w, "El usuario es vendedor o tiene compras asociadas, no es posible eliminarlo.", http.StatusConflict)
// 		return
// 	}
// 	if err != nil {
// 		http.Error(w, "Error al eliminar el usuario: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if filas == 0 {
// 		http.Error(w, "Usuario no encontrado.", http.StatusNotFound)
// 		return
// 	}
// 	w.WriteHeader(http.StatusNoContent)
// }
