package middleware

import (
	"charm.land/log/v2"
	"encoding/json"
	"github.com/nats-io/nats.go"
	db "tan-bit.com/tan-bit/db/sqlc"
	"time"
)

func ExecConsumer() {
	nc, err := nats.Connect("nats://broker:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	// Suscribirse a eventos
	_, err = nc.Subscribe("articulos.events", func(msg *nats.Msg) {
		log.Printf("recibido!")
		var event struct {
			Type     string      `json:"type"`
			Articulo db.Articulo `json:"articulo"`
		}
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			log.Printf("Error decodificando evento: %v", err)
			return
		}

		switch event.Type {
		case "articulo_creado":
			log.Printf("Procesando articulo nuevo: %s", event.Articulo.Nombre)
			// Lógica de negocio aquí...
			time.Sleep(1 * time.Second) // Simular procesamiento
			log.Printf("Articulo %s procesado", event.Articulo.Nombre)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Event consumer running...")
	select {} // Mantener el programa en ejecución
}
