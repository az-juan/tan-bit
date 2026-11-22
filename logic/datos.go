package logic

import (
	"context"
	"database/sql"
	"fmt"

	db "tan-bit.com/tan-bit/db/sqlc"
)

// Datos concentra la interacción con la base usando las queries sqlc.
// Las operaciones de escritura no reciben context ni *sql.DB: usan un
// contexto interno y, cuando hace falta, una transacción propia.
type Datos struct {
	conn *sql.DB
	q    *db.Queries
}

// NuevoDatos crea el acceso a datos a partir del pool de conexiones.
func NuevoDatos(conn *sql.DB) *Datos {
	return &Datos{conn: conn, q: db.New(conn)}
}

// Queries expone las consultas sqlc para las lecturas.
func (d *Datos) Queries() *db.Queries {
	return d.q
}

// ActualizarArticulo aplica la edición solo si el artículo pertenece al vendedor.
func (d *Datos) ActualizarArticulo(idArticulo, idVendedor int32, in ActualizarArticuloInput) (int64, error) {
	return d.q.UpdateArticuloDeVendedor(context.Background(), db.UpdateArticuloDeVendedorParams{
		IDArticulo:  idArticulo,
		IDVendedor:  idVendedor,
		Titulo:      in.Titulo,
		Descripcion: in.Descripcion,
		Condicion:   in.Condicion,
		Categoria:   in.Categoria,
		Precio:      in.Precio,
		PrecioXmr:   in.PrecioXmr,
		Stock:       in.Stock,
		Ubicacion:   in.Ubicacion,
		TipoEntrega: in.TipoEntrega,
		CostoEnvio:  in.CostoEnvio,
		ZonaEntrega: in.ZonaEntrega,
	})
}

// ComprarCarrito registra la compra, descuenta el stock y vacía el carrito
// en una transacción. Devuelve ErrCarritoVacio o un error de stock.
func (d *Datos) ComprarCarrito(idUsuario int32) error {
	ctx := context.Background()
	arts, err := d.q.GetCarritoCompra(ctx, idUsuario)
	if err != nil {
		return err
	}
	if len(arts) == 0 {
		return ErrCarritoVacio
	}
	for _, a := range arts {
		if a.Stock <= 0 {
			return fmt.Errorf("sin stock: %s", a.Titulo)
		}
		if a.Cantidad > a.Stock {
			return fmt.Errorf("stock insuficiente: %s (%d disponibles)", a.Titulo, a.Stock)
		}
	}
	var total float64
	var totalXmr float64
	for _, a := range arts {
		total += a.Precio * float64(a.Cantidad)
		totalXmr += a.PrecioXmr * float64(a.Cantidad)
	}
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op si ya se commiteó

	qtx := db.New(tx)
	compra, err := qtx.CreateCompra(ctx, db.CreateCompraParams{
		IDComprador: idUsuario,
		Total:       total,
		TotalXmr:    totalXmr,
	})
	if err != nil {
		return err
	}
	for _, a := range arts {
		if err := qtx.CreateDetalleCompra(ctx, db.CreateDetalleCompraParams{
			IDCompra:                   compra.IDCompra,
			IDArticulo:                 a.IDArticulo,
			Cantidad:                   a.Cantidad,
			PrecioHistoricoUnitario:    a.Precio,
			PrecioHistoricoUnitarioXmr: a.PrecioXmr,
		}); err != nil {
			return err
		}
		if _, err := qtx.DecrementarStock(ctx, db.DecrementarStockParams{
			IDArticulo: a.IDArticulo,
			Stock:      a.Cantidad,
		}); err != nil {
			return fmt.Errorf("sin stock: %s", a.Titulo)
		}
	}
	if err := qtx.ClearCarrito(ctx, idUsuario); err != nil {
		return err
	}
	return tx.Commit()
}

// ActualizarDatosUsuario aplica el update vía sqlc; conserva el hash actual si
// no se envió una contraseña nueva.
func (d *Datos) ActualizarDatosUsuario(id int32, in ActualizarUsuarioInput) error {
	u, err := d.q.GetUsuarioByID(context.Background(), id)
	if err != nil {
		return err
	}
	password := u.Password
	if in.Password != "" {
		password = HashPassword(in.Password)
	}
	_, err = d.q.UpdateUsuario(context.Background(), db.UpdateUsuarioParams{
		IDUsuario: id,
		Apellido:  in.Apellido,
		Nombre:    in.Nombre,
		Email:     in.Email,
		Password:  password,
	})
	return err
}

// ActualizarVendedor aplica el update vía sqlc.
func (d *Datos) ActualizarVendedor(id int32, in ActualizarVendedorInput) error {
	_, err := d.q.UpdateVendedor(context.Background(), db.UpdateVendedorParams{
		IDVendedor:   id,
		Telefono:     in.Telefono,
		DireccionXmr: in.DireccionXmr,
		Cbu:          in.Cbu,
		Alias:        in.Alias,
	})
	return err
}
