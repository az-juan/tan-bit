package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	sqlc "tan-bit.com/tan-bit/db/sqlc"
)

func main() {
	connStr := "host=database port=5432 user=postgres password=postgres database=tan_bit"
	db, err := sql.Open("pgx", connStr)
	db.Ping()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	queries := sqlc.New(db)
	ctx := context.Background()

	artCreado, err := queries.CreateArticulo(ctx,
		sqlc.CreateArticuloParams{
			Nombre:      "thinkpad t480",
			Precio:      "300.000",
			Descripcion: "notebook thinkpad util para cualquier trabajo de oficio y desarrollo, buena calidad de materiales, sin uso",
			Condicion:   "nuevo",
			Categoria:   "notebook",
			Stock:       1,
			Contacto:    "2494112233",
		})

	if err != nil {
		log.Fatalf("error al crear articulo: %v", err)
	}

	fmt.Printf("Articulo creado: %+v\n", artCreado)

	art, err := queries.GetArticuloByID(ctx, artCreado.ID)
	if err != nil {
		log.Fatalf("error al obtener articulo: %v", err)
	}

	fmt.Printf("Articulo obtenido: %v\n", art)

	arts, err := queries.ListArticulos(ctx)

	if err != nil {
		log.Fatalf("Error al listar articulos: %v", err)
	}
	fmt.Printf("Todos los articulos: %+v\n", arts)

	err = queries.UpdateArticulo(ctx, sqlc.UpdateArticuloParams{
		ID:          artCreado.ID,
		Nombre:      artCreado.Nombre,
		Precio:      "350.000",
		Descripcion: artCreado.Descripcion,
		Condicion:   artCreado.Condicion,
		Categoria:   artCreado.Categoria,
		Stock:       2,
		Contacto:    artCreado.Contacto,
	})

	if err != nil {
		log.Fatalf("Error al actualizar articulo: %v", err)
	}
	fmt.Println("Articulo actualizado sin problemas.")

	artActualizado, err := queries.GetArticuloByID(ctx, artCreado.ID)
	if err != nil {
		log.Fatalf("Error al obtener articulo actualizado: %v", err)
	}
	fmt.Printf("Articulo actualizado: %+v\n", artActualizado)

	err = queries.DeleteArticulo(ctx, artCreado.ID)
	if err != nil {
		log.Fatalf("Error al borrar articulo: %v", err)
	}

	fmt.Println("Articulo borrado sin problemas.")

	_, err = queries.GetArticuloByID(ctx, artCreado.ID)
	if err == sql.ErrNoRows {
		fmt.Println("Articulo no encontrado despues de borrado")
	} else if err != nil {
		log.Fatalf("Error al obtener articulo despues de borrado: %v", err)
	}

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer) //Maneja automaticamente los Content-Type
	port := ":8080"              //de los archivos que sirve
	// http.HandleFunc("/*", handle404)

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
