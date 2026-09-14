package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest" //Paquete std para probar web servers sin abrir puertos reales
	"os"
	sqlc "tan-bit.com/tan-bit/db/sqlc"
	"testing" //Obligatorio en _test.go
)

func TestServidorEstaticoConPG(t *testing.T) { //Firma std en _test.go (Test+Mayus y recibe un puntero)
	if err := godotenv.Load("example.env"); err != nil {
		log.Fatalf("Archivo .env no encontrado")
	}
	connStr := "host=database port=5432 user=" + os.Getenv("PG_USER") + " password=" + os.Getenv("PG_PASSWD") + " database=" + os.Getenv("PG_DB_NAME")
	db, err := sql.Open("pgx", connStr)
	db.Ping()
	if err != nil {
		log.Fatalf("Error al conectarse a la base de datos: %v", err)
	}
	defer db.Close()
	queries := sqlc.New(db)
	ctx := context.Background()

	artCreado, err := queries.CreateArticulo(ctx,
		sqlc.CreateArticuloParams{
			Nombre:      "thinkpad t480",
			Precio:      "300000",
			Descripcion: "notebook thinkpad util para cualquier trabajo de oficio y desarrollo, buena calidad de materiales, sin uso",
			Condicion:   "nuevo",
			Categoria:   "notebook",
			Stock:       1,
			Contacto:    "2494112233",
		})

	if err != nil {
		log.Fatalf("Error al crear articulo: %v", err)
	}

	fmt.Printf("Articulo creado: %+v\n", artCreado)

	art, err := queries.GetArticuloByID(ctx, artCreado.ID)
	if err != nil {
		log.Fatalf("Error al obtener articulo: %v", err)
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

	// Petición HTTP simulada a la raíz "/"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Error al crear la petición: %v", err) //Fatal para arranque
	}

	// Creamos un grabador de respuesta q captura lo que devuelve el handler
	rr := httptest.NewRecorder()

	fileServer.ServeHTTP(rr, req)

	// Verificamos que el código de estado HTTP sea 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Código de estado incorrecto: obtenido %v, esperado %v", status, http.StatusOK)
	}
}
