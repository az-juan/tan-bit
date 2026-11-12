package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tan-bit.com/tan-bit/db/connection"
	sqlc "tan-bit.com/tan-bit/db/sqlc"
	"testing" //Obligatorio en _test.go
)

func TestPG(t *testing.T) { //Firma std en _test.go (Test+Mayus y recibe un puntero)

	db, err := connection.Open()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db) //Instanciamos queries para ejecutarse en nuestro pool
	ctx := context.Background()

	usuarioCreado, err := queries.CreateUsuario(ctx,
		sqlc.CreateUsuarioParams{
			Apellido: "Messi",
			Nombre:   "Leonel",
			Email:    "leo@gmail.com",
		})
	if err != nil {
		log.Fatalf("Error al crear usuario: %v", err)
	}
	fmt.Printf("Usuario creado: %+v\n", usuarioCreado)

	vendedorCreado, err := queries.CreateVendedor(ctx,
		sqlc.CreateVendedorParams{
			IDVendedor:   usuarioCreado.IDUsuario,
			Telefono:     2494777777,
			DireccionXmr: sql.NullString{},
			Cbu:          "24947777770",
			Alias:        "leo_dev_vendedor",
		})
	if err != nil {
		log.Fatalf("Error al crear vendedor: %v", err)
	}
	fmt.Printf("Vendedor creado: %+v\n", vendedorCreado)

	articuloCreado, err := queries.CreateArticulo(ctx,
		sqlc.CreateArticuloParams{
			Titulo:      "Thinkpad T14",
			Descripcion: "notebook empresarial de buena calidad",
			Condicion:   "muy bueno",
			Categoria:   "notebook",
			Precio:      3000000.00,
			Stock:       3,
			IDVendedor:  vendedorCreado.IDVendedor,
			TipoEntrega: "ambos",
		})
	if err != nil {
		log.Fatalf("Error al crear articulo: %v", err)
	}
	fmt.Printf("Articulo creado: %+v\n", articuloCreado)

	art, err := queries.GetArticuloByID(ctx, articuloCreado.IDArticulo)
	if err != nil {
		log.Fatalf("Error al obtener articulo: %v", err)
	}
	fmt.Printf("Articulo obtenido: %v\n", art)

	arts, err := queries.GetArticulos(ctx)
	if err != nil {
		log.Fatalf("Error al listar articulos: %v", err)
	}
	fmt.Printf("Todos los articulos: %+v\n", arts)

	_, err = queries.UpdateArticulo(ctx, sqlc.UpdateArticuloParams{
		IDArticulo:  articuloCreado.IDArticulo,
		Titulo:      articuloCreado.Titulo,
		Descripcion: articuloCreado.Descripcion,
		Condicion:   articuloCreado.Condicion,
		Categoria:   articuloCreado.Categoria,
		Precio:      3500000.00,
		Stock:       2,
		TipoEntrega: "ambos",
	})
	if err != nil {
		log.Fatalf("Error al actualizar articulo: %v", err)
	}
	fmt.Println("Articulo actualizado sin problemas.")

	articuloActualizado, err := queries.GetArticuloByID(ctx, articuloCreado.IDArticulo)
	if err == sql.ErrNoRows {
		fmt.Println("Articulo no encontrado.")
	} else if err != nil {
		log.Fatalf("Error al obtener articulo actualizado: %v", err)
	}
	fmt.Printf("Articulo actualizado: %+v\n", articuloActualizado)

	_, err = queries.DeleteArticulo(ctx, articuloCreado.IDArticulo)
	if err != nil {
		log.Fatalf("Error al borrar articulo: %v", err)
	}
	fmt.Println("Articulo borrado sin problemas.")

	_, err = queries.DeleteVendedor(ctx, vendedorCreado.IDVendedor)
	if err != nil {
		log.Fatalf("Error al borrar vendedor: %v", err)
	}
	fmt.Println("Vendedor borrado sin problemas.")

	_, err = queries.DeleteUsuario(ctx, usuarioCreado.IDUsuario)
	if err != nil {
		log.Fatalf("Error al borrar usuario: %v", err)
	}
	fmt.Println("Usuario borrado sin problemas.")

	_, err = queries.GetArticuloByID(ctx, articuloCreado.IDArticulo)
	if err == sql.ErrNoRows {
		fmt.Println("Articulo no encontrado despues de borrado.")
	} else if err != nil {
		log.Fatalf("Error al obtener articulo despues de borrado: %v", err)
	}

	usuarios, err := queries.GetUsuarios(ctx)
	if err != nil {
		log.Fatalf("Error al obtener usuarios: %v", err)
		return
	}
	fmt.Printf("Usuarios existentes luego del borrado: %v\n", usuarios)
}
