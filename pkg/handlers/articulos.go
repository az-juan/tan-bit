package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tan-bit.com/tan-bit/db/sqlc"
	"time"
)

var articulos = []db.Articulo{
	{1, "Thinkpad T14", "300000.00", "notebook empresarial de buena calidad", "muy bueno", sql.NullString{}, "notebook", 1, "2494001122", sql.NullTime{}},
	{2, "router tplink 300mb", "50000.00", "router ideal para uso domestico", "decente", sql.NullString{}, "internet", 2, "2494987654", sql.NullTime{}},
	{3, "samsung a16", "200000.00", "celular con poco uso", "bueno", sql.NullString{}, "video", 1, "2494123456", sql.NullTime{}},
	{4, "raspberry pi 4", "100000.00", "mini computadora con infinidad de usos", "muy bueno", sql.NullString{}, "electronica", 5, "2494778844", sql.NullTime{}},
	{5, "auriculares logitech", "10000.00", "auriculares de buena calidad, necesitan cambiar almoadillas", "a reparar", sql.NullString{}, "audio", 1, "2494001122", sql.NullTime{}},
}

var nc *nats.Conn

// Ejemplo didáctico: el estado global debe protegerse ante concurrencia
// y reemplazarse por persistencia en una aplicación real.
func ExecBroker() {
	nc, err := nats.Connect("nats://broker:4222")
	if err != nil {
		log.Fatalf("Error en broker: %v", err)
	}
	defer nc.Close()

}

// Manejador para /articulos
func ArtsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getArticulos(w, r)
	case http.MethodPost:
		createArticulo(w, r)
	default:
		http.Error(w, "Metodo no permitido",
			http.StatusMethodNotAllowed)
	}
}

// Manejador para /articulos/{id}
func ArtHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "ID de articulo invalido", http.StatusBadRequest)
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
		http.Error(w, "Metodo no permitido",
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
	var art db.Articulo
	err := json.NewDecoder(r.Body).Decode(&art)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Publicar evento de creación
	type ArticuloCreado struct {
		Type     string      `json:"type"`
		Articulo db.Articulo `json:"articulo"`
		Time     time.Time   `json:"time"`
	}

	event := ArticuloCreado{Type: "articulo_creado", Articulo: art, Time: time.Now()}
	eventData, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "Error creando evento", http.StatusInternalServerError)
		return
	}
	log.Printf("evento creado!")
	if err := nc.Publish("articulos.events", eventData); err != nil {
		http.Error(w, "Error procesando request", http.StatusInternalServerError)
		log.Fatalf("%v", err)
		return
	}
	art.ID = int32(len(articulos)) + 1
	articulos = append(articulos, art)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "procesando",
		"message": "Creacion de articulo en progreso",
	})
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
