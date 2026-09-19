package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"tan-bit.com/tan-bit/db/sqlc"
)

var articulos = []db.Articulo{
	{1, "Thinkpad T14", "300000.00", "notebook empresarial de buena calidad", "muy bueno", sql.NullString{}, "notebook", 1, "2494001122", sql.NullTime{}},
	{2, "router tplink 300mb", "50000.00", "router ideal para uso domestico", "decente", sql.NullString{}, "internet", 2, "2494987654", sql.NullTime{}},
	{3, "samsung a16", "200000.00", "celular con poco uso", "bueno", sql.NullString{}, "video", 1, "2494123456", sql.NullTime{}},
	{4, "raspberry pi 4", "100000.00", "mini computadora con infinidad de usos", "muy bueno", sql.NullString{}, "electronica", 5, "2494778844", sql.NullTime{}},
	{5, "auriculares logitech", "10000.00", "auriculares de buena calidad, necesitan cambiar almoadillas", "a reparar", sql.NullString{}, "audio", 1, "2494001122", sql.NullTime{}},
}

// Ejemplo didáctico: el estado global debe protegerse ante concurrencia
// y reemplazarse por persistencia en una aplicación real.
func main() {
	// Configurar rutas
	http.HandleFunc("/articulos/", artsHandler)
	http.HandleFunc("/articulo/", artHandler)
	// Iniciar servidor
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Manejador para /articulos
func artsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getArticulos(w, r)
	case http.MethodPost:
		createArticulo(w, r)
	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
}

// Manejador para /articulos/{id}
func artHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getArticulo(w, r, id)
	case http.MethodPut:
		updateArticulo(w, r, id)
	case http.MethodDelete:
		deleteArticulo(w, r, id)
	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
}

// GET /articulos - Listar todos los articulo
func getArticulos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articulos)
}

// POST /articulos - Crear nuevo articulo
func createArticulo(w http.ResponseWriter, r *http.Request) {
	var newArticulo db.Articulo
	err := json.NewDecoder(r.Body).Decode(&newArticulo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newArticulo.ID = int32(len(articulos)) + 1
	articulos = append(articulos, newArticulo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newArticulo)
}

// GET /articulos/{id} - Obtener articulo específico
func getArticulo(w http.ResponseWriter, r *http.Request, id int) {
	product, err := findArticuloByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// PUT /articulos/{id} - Actualizar articulo
func updateArticulo(w http.ResponseWriter, r *http.Request, id int) {
	var updatedArticulo db.Articulo
	err := json.NewDecoder(r.Body).Decode(&updatedArticulo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	index, err := findArticuloIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Mantener el ID original
	updatedArticulo.ID = int32(id)
	articulos[index] = updatedArticulo
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedArticulo)
}

// DELETE /articulos/{id} - Eliminar articulo
func deleteArticulo(w http.ResponseWriter, r *http.Request, id int) {
	index, err := findArticuloIndexByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	articulos = append(articulos[:index], articulos[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

// Funciones auxiliares
func findArticuloByID(id int) (*db.Articulo, error) {
	for i := range articulos {
		if articulos[i].ID == int32(id) {
			return &articulos[i], nil
		}
	}
	return nil, errors.New("articulo no encontrado")
}

func findArticuloIndexByID(id int) (int, error) {
	for i, p := range articulos {
		if p.ID == int32(id) {
			return i, nil
		}
	}
	return -1, errors.New("articulo no encontrado")
}
