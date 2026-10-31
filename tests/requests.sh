#! /usr/bin/env bash
set -euo pipefail

echo "Lista de todos los articulos"
curl localhost:8080/articulos
echo

echo "Creacion de nuevo usuario"
curl -X POST http://localhost:8080/usuarios \
-H "Content-Type: application/json" \
-d '{"apellido": "Messi", "nombre": "Leonel", "email": "leo@gmail.com"}'
echo

echo "Hacer al usuario vendedor"
curl -X POST http://localhost:8080/vendedores \
-H "Content-Type: application/json" \
-d '{"id_vendedor": 1, "telefono": 2494777777}'
echo

echo "Listar todos los usuarios"
curl http://localhost:8080/usuarios
echo

echo "Listar todos los vendedores"
curl http://localhost:8080/vendedores
echo

echo "Agregar un nuevo articulo"
curl -X POST http://localhost:8080/articulos \
-H "Content-Type: application/json" \
-d '{"titulo": "Thinkpad T14", "descripcion": "notebook empresarial de buena calidad", "condicion": "muy bueno", "categoria": "notebook", "precio": "300000.00", "stock": 1, "id_vendedor": 1}'
echo

echo "Obtener el articulo 1"
curl localhost:8080/articulos/1
echo

echo "Listar todos los articulos"
curl localhost:8080/articulos
echo