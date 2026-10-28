package tests

import (
	"log"
	"context"
	"database/sql"
	"fmt"
	sqlc "tan-bit.com/tan-bit/db/sqlc"
	"testing" //Obligatorio en _test.go
	"tan-bit.com/tan-bit/db/connection"
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
			"Messi",
			"Leonel",
			"leo@gmail.com",
		})
	if err != nil {
		log.Fatalf("Error al crear usuario: %v", err)
	}
	fmt.Printf("Usuario creado: %+v\n", usuarioCreado)

	vendCreado, err := queries.CreateVendedor(ctx,
		sqlc.CreateVendedorParams{
			1,
			2494777777,
		})
	if err != nil {
		log.Fatalf("Error al crear vendedor: %v", err)
	}
	fmt.Printf("Vendedor creado: %+v\n", vendCreado)

	artCreado, err := queries.CreateArticulo(ctx,
		sqlc.CreateArticuloParams{
			"Thinkpad T14", 
			"notebook empresarial de buena calidad", 
			"muy bueno", 
			"notebook", 
			"300000.00", 
			1, 
			1,
		})
	if err != nil {
		log.Fatalf("Error al crear articulo: %v", err)
	}
	fmt.Printf("Articulo creado: %+v\n", artCreado)

	art, err := queries.GetArticuloByID(ctx, 1)
	if err != nil {
		log.Fatalf("Error al obtener articulo: %v", err)
	}
	fmt.Printf("Articulo obtenido: %v\n", art)

	arts, err := queries.GetArticulos(ctx)
	if err != nil {
		log.Fatalf("Error al listar articulos: %v", err)
	}
	fmt.Printf("Todos los articulos: %+v\n", arts)

	err = queries.UpdateArticulo(ctx, sqlc.UpdateArticuloParams{
		IDArticulo:  1,
		Titulo:      artCreado.Titulo,
		Descripcion: artCreado.Descripcion,
		Condicion:   artCreado.Condicion,
		Precio:      "350.000",
		Stock:       2,
	})
	if err != nil {
		log.Fatalf("Error al actualizar articulo: %v", err)
	}
	fmt.Println("Articulo actualizado sin problemas.")

	artActualizado, err := queries.GetArticuloByID(ctx, 1)
	if err != nil {
		log.Fatalf("Error al obtener articulo actualizado: %v", err)
	}
	fmt.Printf("Articulo actualizado: %+v\n", artActualizado)

	err = queries.DeleteArticulo(ctx, 1)
	if err != nil {
		log.Fatalf("Error al borrar articulo: %v", err)
	}
	fmt.Println("Articulo borrado sin problemas.")

	err = queries.DeleteVendedor(ctx, vendCreado.IDVendedor)
	if err != nil {
		log.Fatalf("Error al borrar vendedor: %v", err)
	}
	fmt.Println("Vendedor borrado sin problemas.")

	err = queries.DeleteUsuario(ctx, 1)
	if err != nil {
		log.Fatalf("Error al borrar usuario: %v", err)
	}
	fmt.Println("Usuario borrado sin problemas.")

	_, err = queries.GetArticuloByID(ctx, 1)
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
	fmt.Printf("Usuarios existentes: %v\n", usuarios)
}