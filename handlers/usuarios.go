package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	db "tan-bit.com/tan-bit/db/sqlc"
	"tan-bit.com/tan-bit/logic"
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
		views.Registrarse("", "usuario").Render(r.Context(), w)
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

	role := r.FormValue("role")
	if role != "usuario" && role != "vendedor" {
		role = "usuario"
	}

	terminos := r.FormValue("terminos")
	if terminos != "on" {
		views.Pagina(r, views.Registrarse("Debes aceptar los términos del servicio y la política de privacidad.", role)).Render(r.Context(), w)
		return
	}

	if role == "vendedor" {
		if r.FormValue("telefono") == "" || r.FormValue("cbu") == "" || r.FormValue("alias") == "" {
			views.Pagina(r, views.Registrarse("Para registrarte como vendedor necesitás completar TELÉFONO, CBU y ALIAS.", role)).Render(r.Context(), w)
			return
		}
	}

	params := db.CreateUsuarioParams{
		Apellido: r.FormValue("apellido"),
		Nombre:   r.FormValue("nombre"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	nuevoUsuario, err := h.queries.CreateUsuario(r.Context(), params)
	if err != nil {
		if logic.EsViolacionUnique(err) {
			col := logic.ColumnaViolacionUnique(err)
			var msg string
			switch col {
			case "cbu":
				msg = "Este CBU ya está asociado a otra cuenta vendedora."
			case "alias":
				msg = "Este alias ya está en uso."
			case "telefono":
				msg = "Este teléfono ya está asociado a otra cuenta vendedora."
			default:
				msg = "El email ya está registrado. ¿Tenés cuenta? Iniciá sesión."
			}
			views.Registrarse(msg, role).Render(r.Context(), w)
			return
		}
		http.Error(w, "Error al crear el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if role == "vendedor" {
		telefono, err := strconv.ParseInt(r.FormValue("telefono"), 10, 64)
		if err != nil {
			views.Registrarse("TELÉFONO debe ser un número válido.", role).Render(r.Context(), w)
			return
		}
		cbu := r.FormValue("cbu")
		direccionXmr := sql.NullString{}
		if v := r.FormValue("direccion_xmr"); v != "" {
			direccionXmr = sql.NullString{String: v, Valid: true}
		}
		_, err = h.queries.CreateVendedor(r.Context(), db.CreateVendedorParams{
			IDVendedor:   nuevoUsuario.IDUsuario,
			Telefono:     telefono,
			DireccionXmr: direccionXmr,
			Cbu:          cbu,
			Alias:        r.FormValue("alias"),
		})
		if err != nil {
			http.Error(w, "Error al crear el vendedor: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{Name: "user_id", Value: strconv.Itoa(int(nuevoUsuario.IDUsuario)), Path: "/"})
	// http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
}

func (h *UsuarioHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		views.Pagina(r, views.Login()).Render(r.Context(), w)
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
	// views.NotificationList(logic.GetNotifications()).Render(r.Context(), w)
}

func (h *UsuarioHandler) MiCuentaHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	nombre, apellido, email := "", "", ""
	anonTelefono := ""
	alias := ""
	cbu := ""
	direccionXmr := ""
	anarticulos := []db.GetArticulosByVendedorRow{}
	anfavoritos := []db.GetFavoritosByUsuarioRow{}
	ancompras := []db.GetComprasByUsuarioRow{}

	if u, err := h.queries.GetUsuarioByID(r.Context(), int32(id)); err == nil {
		nombre = u.Nombre
		apellido = u.Apellido
		email = u.Email
	}
	if v, err := h.queries.GetVendedorByID(r.Context(), int32(id)); err == nil {
		anonTelefono = strconv.FormatInt(v.Telefono, 10)
		alias = v.Alias
		cbu = v.Cbu
		if v.DireccionXmr.Valid {
			direccionXmr = v.DireccionXmr.String
		}
	}
	if a, err := h.queries.GetArticulosByVendedor(r.Context(), int32(id)); err == nil {
		anarticulos = a
	}
	if f, err := h.queries.GetFavoritosByUsuario(r.Context(), int32(id)); err == nil {
		anfavoritos = f
	}
	if c, err := h.queries.GetComprasByUsuario(r.Context(), int32(id)); err == nil {
		ancompras = c
	}
	views.Pagina(r, views.MiCuenta(nombre, apellido, email, anonTelefono, alias, cbu, direccionXmr, anarticulos, anfavoritos, ancompras)).Render(r.Context(), w)
}
