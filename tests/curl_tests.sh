#! /usr/bin/env bash
set -euo pipefail

curl localhost:8080/articulos
echo
curl localhost:8080/articulos/2
echo
curl -i -X DELETE localhost:8080/articulos/2
echo
curl -X POST http://localhost:8080/articulos \
-H "Content-Type: application/json" \
-d '{"nombre": "LoRa RNode", "precio": "50000.00", "descripcion": "dispositivo para comunicacion offgrid", "condicion": "nuevo", "categoria": "electronica", "stock": 1, "contacto": "offgrid_tan@gmail.com"}'
