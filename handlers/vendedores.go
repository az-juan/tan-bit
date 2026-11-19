package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	db "tan-bit.com/tan-bit/db/sqlc"
	"tan-bit.com/tan-bit/logic"
	"tan-bit.com/tan-bit/views"
)

// ajustar string para la ruta de la imagen
func sanitize(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			b.WriteRune('_')
		}
	}
	return b.String()
}

type VendedorHandler struct {
	queries *db.Queries
}

func NewVendedorHandler(queries *db.Queries) *VendedorHandler {
	return &VendedorHandler{
		queries: queries,
	}
}

// ServeHTTP para vendedor handler
func (h *VendedorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/vender":
		h.VenderHandler(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

// Manejador de página de venta
func (h *VendedorHandler) VenderHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := logic.GetCookieOrRedirect(w, r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	v, err := h.queries.GetVendedorByID(r.Context(), int32(id))
	if err != nil || v.Telefono == 0 {
		return
	}
	if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			r.ParseForm()
		}
		data := views.ParseVenderData(r)

		// Campos obligatorios
		var missing []string
		if data.Titulo == "" {
			missing = append(missing, "título")
		}
		if data.Descripcion == "" {
			missing = append(missing, "descripción")
		}
		if data.Precio == "" {
			data.Precio = "0"
		}
		if data.PrecioXMR == "" {
			data.PrecioXMR = "0"
		}
		if data.Ubicacion == "" {
			missing = append(missing, "ubicación")
		}
		if data.ZonaEntrega == "" {
			missing = append(missing, "zona de entrega")
		}
		if data.TipoEntrega == "" {
			missing = append(missing, "tipo de entrega")
		}
		if len(missing) > 0 {
			errMsg := "Faltan datos obligatorios: " + strings.Join(missing, ", ")
			views.Pagina(r, views.Vender(data, views.CategoriasDB(), views.CondicionesDB(), errMsg)).Render(r.Context(), w)
			return
		}
		idVendedor := int32(0)
		if cookie, err := r.Cookie("user_id"); err == nil {
			if v, err := strconv.Atoi(cookie.Value); err == nil {
				idVendedor = int32(v)
			}
		}
		if idVendedor == 0 {
			http.Redirect(w, r, "/mi_cuenta", http.StatusSeeOther)
			return
		}

		// Parseo de valores
		precio := 0.0
		if v := r.FormValue("precio"); v != "" {
			precio, _ = strconv.ParseFloat(v, 64)
		}
		precioXmr := 0.0
		if v := r.FormValue("precio_xmr"); v != "" {
			precioXmr, _ = strconv.ParseFloat(v, 64)
		}
		costoEnvio := 0.0
		if v := data.CostoEnvio; v != "" {
			costoEnvio, _ = strconv.ParseFloat(v, 64)
		}

		// guardar imagen
		if f, fh, err := r.FormFile("fotos"); err == nil {
			defer f.Close()
			ext := filepath.Ext(fh.Filename)
			if ext == "" {
				ext = ".jpg"
			}
			os.MkdirAll("uploads", 0755)
			nombre := fmt.Sprintf("%d_%s_1%s", idVendedor, sanitize(data.Titulo), ext)
			dst := filepath.Join("uploads", nombre)
			out, err := os.Create(dst)
			if err == nil {
				io.Copy(out, f)
				out.Close()
				data.RutaImagen = "/uploads/" + nombre
			}
		}

		stock := int32(1)
		if v, err := strconv.Atoi(data.Stock); err == nil {
			stock = int32(v)
		}
		params := db.CreateArticuloParams{
			Titulo:      data.Titulo,
			Descripcion: data.Descripcion,
			Condicion:   data.Condicion,
			Categoria:   data.Categoria,
			Precio:      precio,
			PrecioXmr:   precioXmr,
			Stock:       stock,
			IDVendedor:  idVendedor,
			Ubicacion:   data.Ubicacion,
			TipoEntrega: data.TipoEntrega,
			CostoEnvio:  sql.NullFloat64{Float64: costoEnvio, Valid: costoEnvio != 0},
			ZonaEntrega: sql.NullString{String: data.ZonaEntrega, Valid: data.ZonaEntrega != ""},
			RutaImagen:  sql.NullString{String: data.RutaImagen, Valid: data.RutaImagen != ""},
		}
		if _, err := h.queries.CreateArticulo(r.Context(), params); err != nil {
			errMsg := "No se pudo publicar el artículo: " + err.Error()
			views.Pagina(r, views.Vender(data, views.CategoriasDB(), views.CondicionesDB(), errMsg)).Render(r.Context(), w)
			return
		}
		// http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	views.Pagina(r, views.Vender(views.VenderData{}, views.CategoriasDB(), views.CondicionesDB(), "")).Render(r.Context(), w)
}
