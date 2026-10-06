#! /usr/bin/env bash
set -euo pipefail

gum style --foreground 150 --bold "Listar todos los articulos"
curl localhost:8080/articulos
echo

gum style --foreground 150 --bold "Obtener el articulo 2"
curl localhost:8080/articulos/2
echo

gum style --foreground 150 --bold "Borrar el articulo 2"
curl -i -X DELETE localhost:8080/articulos/2
echo

gum style --foreground 150 --bold "Agregar un nuevo articulo"
curl -X POST http://localhost:8080/articulos \
-H "Content-Type: application/json" \
-d '{"nombre": "LoRa RNode", "precio": "50000.00", "descripcion": "dispositivo para comunicacion offgrid", "condicion": "nuevo", "categoria": "electronica", "stock": 1, "contacto": "offgrid_tan@gmail.com"}'
echo
