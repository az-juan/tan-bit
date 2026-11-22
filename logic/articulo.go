package logic

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	db "tan-bit.com/tan-bit/db/sqlc"
)

// ErrCarritoVacio indica que no hay piezas para comprar.
var ErrCarritoVacio = errors.New("el carrito está vacío")

// CategoriaValida verifica que la categoría exista en el enum de la DB.
func CategoriaValida(s string) bool {
	switch db.Categoria(s) {
	case db.CategoriaNotebook, db.CategoriaPantalla, db.CategoriaElectronica,
		db.CategoriaPeriferico, db.CategoriaAccesorio, db.CategoriaAudio,
		db.CategoriaVideo, db.CategoriaGaming, db.CategoriaComunicacion,
		db.CategoriaAnalogico, db.CategoriaDecoracion:
		return true
	}
	return false
}

// CondicionValida verifica que la condición exista en el enum de la DB.
func CondicionValida(s string) bool {
	switch db.Condicion(s) {
	case db.CondicionNuevo, db.CondicionMuybueno, db.CondicionUsado,
		db.CondicionReacondicionado, db.CondicionAreparar:
		return true
	}
	return false
}

// EntregaValida verifica que el tipo de entrega exista en el enum de la DB.
func EntregaValida(s string) bool {
	switch db.Entrega(s) {
	case db.EntregaEnvioadomicilio, db.EntregaRetiroenpersona, db.EntregaAmbos:
		return true
	}
	return false
}

// ActualizarArticuloInput son los campos editables de una publicación.
type ActualizarArticuloInput struct {
	Titulo      string
	Descripcion string
	Condicion   db.Condicion
	Categoria   db.Categoria
	Precio      float64
	PrecioXmr   float64
	Stock       int32
	Ubicacion   string
	TipoEntrega db.Entrega
	CostoEnvio  sql.NullFloat64
	ZonaEntrega sql.NullString
}

// ParseActualizarArticulo valida el formulario de edición de un artículo.
func ParseActualizarArticulo(r *http.Request) (ActualizarArticuloInput, bool) {
	in := ActualizarArticuloInput{}
	in.Titulo = strings.TrimSpace(r.FormValue("titulo"))
	if in.Titulo == "" {
		return in, false
	}
	in.Descripcion = strings.TrimSpace(r.FormValue("descripcion"))
	if in.Descripcion == "" {
		return in, false
	}
	if !CategoriaValida(r.FormValue("categoria")) {
		return in, false
	}
	if !CondicionValida(r.FormValue("condicion")) {
		return in, false
	}
	if !EntregaValida(r.FormValue("tipo_entrega")) {
		return in, false
	}
	precio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
	if err != nil || precio <= 0 {
		return in, false
	}
	precioXmr, err := strconv.ParseFloat(r.FormValue("precio_xmr"), 64)
	if err != nil || precioXmr < 0 {
		return in, false
	}
	stock, err := strconv.ParseInt(r.FormValue("stock"), 10, 32)
	if err != nil || stock < 0 {
		return in, false
	}
	in.Ubicacion = strings.TrimSpace(r.FormValue("ubicacion"))
	if in.Ubicacion == "" {
		return in, false
	}
	costo, _ := strconv.ParseFloat(r.FormValue("costo_envio"), 64)
	in.CostoEnvio = sql.NullFloat64{Float64: costo, Valid: costo > 0}
	zona := strings.TrimSpace(r.FormValue("zona_entrega"))
	in.ZonaEntrega = sql.NullString{String: zona, Valid: zona != ""}
	in.Condicion = db.Condicion(r.FormValue("condicion"))
	in.Categoria = db.Categoria(r.FormValue("categoria"))
	in.Precio = precio
	in.PrecioXmr = precioXmr
	in.Stock = int32(stock)
	in.TipoEntrega = db.Entrega(r.FormValue("tipo_entrega"))
	return in, true
}
