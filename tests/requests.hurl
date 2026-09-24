# requests.hurl

# 1. Obtener todos los articulos
GET http://localhost:8080/articulos
HTTP 200
[Asserts]
# jsonpath "$.length()" == 5
jsonpath "$[0].nombre" == "Thinkpad T14"

# 2. Obtener el articulo con id 2
GET http://localhost:8080/articulos/2
HTTP 200
[Asserts]
jsonpath "$.nombre" == "router tplink 300mb"

# 3. Eliminar el articulo con id 2
DELETE http://localhost:8080/articulos/2
HTTP 204

# 4. Crear un nuevo articulo
POST http://localhost:8080/articulos
Content-Type: application/json
{
  "nombre": "LoRa RNode",
  "precio": 50000.00,
  "descripcion": "dispositivo para comunicacion offgrid",
  "condicion": "nuevo",
  "categoria": "comunicacion",
  "stock": 1,
  "contacto": "offgrid_tan@gmail.com"
}
HTTP 201
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.nombre" == "LoRa RNode"
