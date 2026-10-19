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
-d '{"titulo": "Thinkpad T14", "descripcion": "notebook empresarial de buena calidad", "condicion": "muy bueno", "categoria": "notebook", "precio": "300000.00", "stock": 1, "id_vendedor": 1}'
echo
