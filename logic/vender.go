package logic

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	db "tan-bit.com/tan-bit/db/sqlc"
)

// VenderInput es la estructura de datos para el panel de vendedor
type VenderInput struct {
	// IDVendedor int32
	// Nombre       string
	// Telefono     string
	// Alias        string
	// Cbu          string
	// DireccionXmr sql.NullString
	IDArticulo  int32
	Titulo      string
	Categoria   string
	Condicion   string
	Descripcion string
	RutaImagen  string
	Ubicacion   string
	Precio      float64
	PrecioXmr   float64
	Stock       int32
	TipoEntrega string
	CostoEnvio  float64
	ZonaEntrega string
}

func ParseVender(r *http.Request, id_vendedor int) (VenderInput, error) {
	var missing []string

	titulo := r.FormValue("titulo")
	if titulo == "" {
		missing = append(missing, "título")
	}

	categoria := r.FormValue("categoria")

	condicion := r.FormValue("condicion")

	descripcion := r.FormValue("descripcion")
	if descripcion == "" {
		missing = append(missing, "descripción")
	}

	ubicacion := r.FormValue("ubicacion")
	if ubicacion == "" {
		missing = append(missing, "ubicación")
	}

	precio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
	if err != nil || precio <= 0 {
		missing = append(missing, "precio")
	}

	precio_xmr, err := strconv.ParseFloat(r.FormValue("precio_xmr"), 64)
	if err != nil || precio_xmr < 0 {
		precio_xmr = 0
	}

	stock, err := strconv.ParseInt(r.FormValue("stock"), 0, 32)
	if err != nil || stock <= 0 {
		missing = append(missing, "stock")
	}

	tipo_entrega := r.FormValue("tipo_entrega")
	if tipo_entrega == "" {
		missing = append(missing, "tipo de entrega")
	}

	costo_envio, err := strconv.ParseFloat(r.FormValue("costo_envio"), 64)
	if err != nil || costo_envio < 0 {
		costo_envio = 0
	}

	zona_entrega := r.FormValue("zona_entrega")
	if zona_entrega == "" {
		missing = append(missing, "zona de entrega")
	}

	if len(missing) > 0 {
		return VenderInput{}, err
	}

	// guardar imagen
	ruta_imagen := ""
	if f, fh, err := r.FormFile("fotos"); err == nil {
		defer f.Close()
		ext := filepath.Ext(fh.Filename)
		if ext == "" {
			ext = ".jpg"
		}
		os.MkdirAll("uploads", 0755)
		nombre := fmt.Sprintf("%d_%s_1%s", id_vendedor, sanitize(titulo), ext)
		dst := filepath.Join("uploads", nombre)
		out, err := os.Create(dst)
		if err == nil {
			io.Copy(out, f)
			out.Close()
			ruta_imagen = "/uploads/" + nombre
		}
	}

	return VenderInput{
		// IDVendedor:  id_vendedor,
		Titulo:      titulo,
		Descripcion: descripcion,
		RutaImagen:  ruta_imagen,
		Condicion:   condicion,
		Categoria:   categoria,
		Precio:      precio,
		PrecioXmr:   precio_xmr,
		Stock:       int32(stock),
		Ubicacion:   ubicacion,
		TipoEntrega: tipo_entrega,
		CostoEnvio:  costo_envio,
		ZonaEntrega: zona_entrega,
	}, nil
}

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

// CategoriasDB devuelve la lista de categorías disponibles en la base de datos
func CategoriasDB() []db.Categoria {
	return []db.Categoria{
		db.CategoriaNotebook,
		db.CategoriaPantalla,
		db.CategoriaElectronica,
		db.CategoriaPeriferico,
		db.CategoriaAccesorio,
		db.CategoriaAudio,
		db.CategoriaVideo,
		db.CategoriaGaming,
		db.CategoriaComunicacion,
		db.CategoriaAnalogico,
		db.CategoriaDecoracion,
	}
}

// CondicionesDB devuelve la lista de condiciones disponibles en la base de datos
func CondicionesDB() []db.Condicion {
	return []db.Condicion{
		db.CondicionNuevo,
		db.CondicionMuybueno,
		db.CondicionUsado,
		db.CondicionReacondicionado,
		db.CondicionAreparar,
	}
}
