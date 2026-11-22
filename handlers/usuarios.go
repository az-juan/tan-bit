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
	datos *logic.Datos
}

func NewUsuarioHandler(datos *logic.Datos) *UsuarioHandler {
	return &UsuarioHandler{datos: datos}
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
	case "/logout":
		h.LogoutHandler(w, r)
		return
	case "/mi_cuenta/delete":
		h.DeleteAccountHandler(w, r)
		return
	default:
		NotFoundHandler(w, r)
	}
}

// notificar encola un aviso para la sesión actual.
func notificar(w http.ResponseWriter, r *http.Request, tipo logic.NotificationType, label, titulo, mensaje, href, texto string) {
	logic.AddNotification(logic.SessionID(w, r), logicaNotificacion(tipo, label, titulo, mensaje, href, texto))
}

// notificarUsuario encola un aviso atado al id de usuario recién creado.
// Se usa cuando la cookie user_id se escribe en la misma respuesta.
func notificarUsuario(id int32, tipo logic.NotificationType, label, titulo, mensaje, href, texto string) {
	logic.AddNotification("user:"+strconv.Itoa(int(id)), logicaNotificacion(tipo, label, titulo, mensaje, href, texto))
}

func logicaNotificacion(tipo logic.NotificationType, label, titulo, mensaje, href, texto string) logic.Notification {
	return logic.Notification{
		Type:       tipo,
		Label:      label,
		Title:      titulo,
		Message:    mensaje,
		ActionHref: href,
		ActionText: texto,
	}
}

// setSessionCookie deja la cookie de sesión con atributos seguros.
func setSessionCookie(w http.ResponseWriter, id int32) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(int(id)),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// clearSessionCookie invalida la cookie de sesión.
func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *UsuarioHandler) RegistrarseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		views.Pagina(w, r, views.Registrarse()).Render(r.Context(), w)
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			notificar(w, r, logic.NotificationError, "error", "Formulario inválido", "No pudimos procesar el envío.", "/registrarse", "Reintentar")
			http.Redirect(w, r, "/registrarse", http.StatusSeeOther)
			return
		}
		data := logic.ParseRegistrarse(r)
		if data == nil {
			notificar(w, r, logic.NotificationError, "error", "No pudimos registrarte", "Revisá los datos: email válido, contraseñas iguales (mínimo 6), nombre y apellido. Si elegís vendedor, también teléfono, CBU de 13 dígitos y alias.", "/registrarse", "Reintentar")
			http.Redirect(w, r, "/registrarse", http.StatusSeeOther)
			return
		}
		params := db.CreateUsuarioParams{
			Apellido: data.Apellido,
			Nombre:   data.Nombre,
			Email:    data.Email,
			Password: logic.HashPassword(data.Password),
		}
		usuario, err := h.datos.Queries().CreateUsuario(r.Context(), params)
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
					msg = "El email ya está registrado."
				}
				notificar(w, r, logic.NotificationError, "error", "No pudimos registrarte", msg, "/login", "Iniciar sesión")
				http.Redirect(w, r, "/registrarse", http.StatusSeeOther)
				return
			}
			notificar(w, r, logic.NotificationError, "error", "Error interno", "No pudimos crear la cuenta. Intentá de nuevo más tarde.", "/registrarse", "Reintentar")
			http.Error(w, "Error al crear el usuario: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rol := "usuario"
		if data.Role == "vendedor" && data.Telefono != "" && data.CBU != "" && data.Alias != "" {
			rol = "vendedor"
			telefono, err := strconv.ParseInt(data.Telefono, 10, 64)
			if err != nil {
				notificar(w, r, logic.NotificationError, "error", "Teléfono inválido", "El teléfono solo puede contener dígitos.", "/registrarse", "Reintentar")
				http.Redirect(w, r, "/registrarse", http.StatusSeeOther)
				return
			}
			_, err = h.datos.Queries().CreateVendedor(r.Context(), db.CreateVendedorParams{
				IDVendedor:   usuario.IDUsuario,
				Telefono:     telefono,
				DireccionXmr: sql.NullString{String: data.DireccionXmr, Valid: data.DireccionXmr != ""},
				Cbu:          data.CBU,
				Alias:        data.Alias,
			})
			if err != nil {
				msg := "La cuenta fue creada, pero no pudimos completar los datos de vendedor."
				if logic.EsViolacionUnique(err) {
					msg = mensajeViolacion(logic.ColumnaViolacionUnique(err))
				}
				logic.SetRolCookie(w, "usuario")
				setSessionCookie(w, usuario.IDUsuario)
				notificar(w, r, logic.NotificationError, "error", "Datos de vendedor ya en uso", msg, "/mi_cuenta", "Revisar mis datos")
				http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
				return
			}
		}
		logic.SetRolCookie(w, rol)
		setSessionCookie(w, usuario.IDUsuario)
		notificarUsuario(usuario.IDUsuario, logic.NotificationSuccess, "éxito", "Cuenta creada", "Ya podés comprar, guardar favoritos y publicar piezas.", "/mi_cuenta", "Ir a mi cuenta")
		if logic.IsHtmxRequest(r) {
			w.Header().Set("HX-Redirect", "/mi_cuenta")
			return
		}
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
}

func (h *UsuarioHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		views.Pagina(w, r, views.Login()).Render(r.Context(), w)
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			notificar(w, r, logic.NotificationError, "error", "Formulario inválido", "No pudimos procesar el envío.", "/login", "Reintentar")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		u, err := h.datos.Queries().GetUsuarioByEmail(r.Context(), email)
		if err != nil {
			notificar(w, r, logic.NotificationError, "error", "Ese email no está registrado", "Revisá la dirección o creá una cuenta.", "/registrarse", "Crear cuenta")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if !logic.CheckPassword(u.Password, password) {
			notificar(w, r, logic.NotificationError, "error", "Contraseña incorrecta", "La contraseña no coincide con esa cuenta.", "/login", "Reintentar")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		rol := "usuario"
		if _, err := h.datos.Queries().GetVendedorByID(r.Context(), u.IDUsuario); err == nil {
			rol = "vendedor"
		}
		logic.SetRolCookie(w, rol)
		setSessionCookie(w, u.IDUsuario)
		notificarUsuario(u.IDUsuario, logic.NotificationInfo, "info", "Sesión iniciada", "Bienvenido de nuevo.", "/mi_cuenta", "Ir a mi cuenta")
		if logic.IsHtmxRequest(r) {
			w.Header().Set("HX-Redirect", "/mi_cuenta")
			return
		}
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
}

func (h *UsuarioHandler) MiCuentaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.actualizarMiCuenta(w, r)
		return
	}
	h.miCuentaGET(w, r)
}

// actualizarMiCuenta enruta el POST entre datos personales y cambio de contraseña.
func (h *UsuarioHandler) actualizarMiCuenta(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		notificar(w, r, logic.NotificationError, "error", "Formulario inválido", "No pudimos procesar el envío.", "/mi_cuenta", "Reintentar")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	if r.FormValue("password_actual") != "" || r.FormValue("password_nueva") != "" {
		h.actualizarContrasenia(w, r, id)
		return
	}
	h.actualizarDatos(w, r, id)
}

func (h *UsuarioHandler) miCuentaGET(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	usuario := views.UsuarioData{}
	if u, err := h.datos.Queries().GetUsuarioByID(r.Context(), int32(id)); err == nil {
		usuario.Nombre = u.Nombre
		usuario.Apellido = u.Apellido
		usuario.Email = u.Email
	}
	if v, err := h.datos.Queries().GetVendedorByID(r.Context(), int32(id)); err == nil {
		usuario.Telefono = v.Telefono
		usuario.Alias = v.Alias
		usuario.Cbu = v.Cbu
		usuario.DireccionXmr = v.DireccionXmr
	}
	a, err := h.datos.Queries().GetArticulosByVendedor(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al obtener los articulos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data_a := make([]views.ArticuloData, len(a))
	for i, a := range a {
		data_a[i] = views.ArticuloData{
			IDArticulo:       a.IDArticulo,
			Titulo:           a.Titulo,
			Descripcion:      a.Descripcion,
			Condicion:        string(a.Condicion),
			Categoria:        string(a.Categoria),
			RutaImagen:       a.RutaImagen,
			Precio:           a.Precio,
			PrecioXmr:        a.PrecioXmr,
			Stock:            a.Stock,
			FechaPublicacion: a.FechaPublicacion,
			Ubicacion:        a.Ubicacion,
			TipoEntrega:      string(a.TipoEntrega),
			CostoEnvio:       a.CostoEnvio,
			ZonaEntrega:      a.ZonaEntrega,
		}
	}
	f, err := h.datos.Queries().GetFavoritosByUsuario(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al obtener los favoritos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data_f := make([]views.ArticuloData, len(f))
	for i, a := range f {
		data_f[i] = views.ArticuloData{
			IDArticulo:       a.IDArticulo,
			Titulo:           a.Titulo,
			Descripcion:      a.Descripcion,
			Condicion:        string(a.Condicion),
			Categoria:        string(a.Categoria),
			RutaImagen:       a.RutaImagen,
			Precio:           a.Precio,
			PrecioXmr:        a.PrecioXmr,
			Stock:            a.Stock,
			FechaPublicacion: a.FechaPublicacion,
			Ubicacion:        a.Ubicacion,
			TipoEntrega:      string(a.TipoEntrega),
			CostoEnvio:       a.CostoEnvio,
			ZonaEntrega:      a.ZonaEntrega,
		}
	}
	c, err := h.datos.Queries().GetComprasByUsuario(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al obtener las compras: "+err.Error(), http.StatusInternalServerError)
		return
	}
	views.Pagina(w, r, views.MiCuenta(usuario, data_a, data_f, c)).Render(r.Context(), w)
}

// mensajeViolacion describe la columna que violó una restricción UNIQUE.
func mensajeViolacion(col string) string {
	switch col {
	case "cbu":
		return "Este CBU ya está asociado a otra cuenta vendedora."
	case "alias":
		return "Este alias ya está en uso."
	case "telefono":
		return "Este teléfono ya está asociado a otra cuenta vendedora."
	default:
		return "El email ya está registrado."
	}
}

// actualizarDatos guarda los datos personales y, si corresponde, los de vendedor.
func (h *UsuarioHandler) actualizarDatos(w http.ResponseWriter, r *http.Request, id int) {
	in, ok := logic.ParseActualizarUsuario(r)
	if !ok {
		notificar(w, r, logic.NotificationError, "error", "No pudimos guardar los datos", "Revisá nombre, apellido y email.", "/mi_cuenta", "Reintentar")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	if _, err := h.datos.Queries().GetVendedorByID(r.Context(), int32(id)); err == nil {
		vin, vok := logic.ParseActualizarVendedor(r)
		if !vok {
			notificar(w, r, logic.NotificationError, "error", "No pudimos guardar los datos de vendedor", "Revisá teléfono, CBU de 13 dígitos y alias.", "/mi_cuenta", "Reintentar")
			http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
			return
		}
		if err := h.datos.ActualizarVendedor(int32(id), vin); err == nil {
			logic.SetRolCookie(w, "vendedor")
		} else {
			msg := "Ocurrió un error al guardar los datos de vendedor."
			if logic.EsViolacionUnique(err) {
				msg = mensajeViolacion(logic.ColumnaViolacionUnique(err))
			}
			notificar(w, r, logic.NotificationError, "error", "No pudimos guardar los datos de vendedor", msg, "/mi_cuenta", "Reintentar")
			http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
			return
		}
	}
	if err := h.datos.ActualizarDatosUsuario(int32(id), in); err != nil {
		msg := "Ocurrió un error al guardar tus datos."
		if logic.EsViolacionUnique(err) {
			msg = mensajeViolacion(logic.ColumnaViolacionUnique(err))
		}
		notificar(w, r, logic.NotificationError, "error", "No pudimos guardar los datos", msg, "/mi_cuenta", "Reintentar")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	notificar(w, r, logic.NotificationSuccess, "éxito", "Datos actualizados", "Tu perfil quedó al día.", "/mi_cuenta", "Ver mi cuenta")
	http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
}

// actualizarContrasenia verifica la contraseña actual y guarda el hash nuevo.
func (h *UsuarioHandler) actualizarContrasenia(w http.ResponseWriter, r *http.Request, id int) {
	u, err := h.datos.Queries().GetUsuarioByID(r.Context(), int32(id))
	if err != nil || !logic.CheckPassword(u.Password, r.FormValue("password_actual")) {
		notificar(w, r, logic.NotificationError, "error", "Contraseña actual incorrecta", "Verificá tu contraseña actual.", "/mi_cuenta", "Reintentar")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	nueva := r.FormValue("password_nueva")
	if len(nueva) < 6 {
		notificar(w, r, logic.NotificationError, "error", "Contraseña demasiado corta", "La nueva contraseña debe tener al menos 6 caracteres.", "/mi_cuenta", "Reintentar")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	if err := h.datos.ActualizarDatosUsuario(int32(id), logic.ActualizarUsuarioInput{
		Nombre: u.Nombre, Apellido: u.Apellido, Email: u.Email, Password: nueva,
	}); err != nil {
		notificar(w, r, logic.NotificationError, "error", "No pudimos cambiar la contraseña", "Intentá de nuevo en unos segundos.", "/mi_cuenta", "Reintentar")
		http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
		return
	}
	notificar(w, r, logic.NotificationSuccess, "éxito", "Contraseña actualizada", "Tu contraseña se cambió correctamente.", "/mi_cuenta", "Ver mi cuenta")
	http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
}

func (h *UsuarioHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	notificar(w, r, logic.NotificationInfo, "info", "Sesión cerrada", "Te esperamos de vuelta.", "/login", "Iniciar sesión")
	logic.ClearRolCookie(w)
	logic.ClearRolCookie(w)
	clearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *UsuarioHandler) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		return
	}
	if _, err := h.datos.Queries().DeleteUsuario(r.Context(), int32(id)); err != nil {
		http.Error(w, "Error al borrar la cuenta: "+err.Error(), http.StatusInternalServerError)
		return
	}
	notificar(w, r, logic.NotificationInfo, "info", "Cuenta eliminada", "Tu cuenta y tus datos fueron borrados.", "/", "Volver al inicio")
	clearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
