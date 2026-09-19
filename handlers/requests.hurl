# requests.hurl

# 1. Obtener todos los articulos
GET http://localhost:8080/articulos
HTTP 200
[Asserts]
jsonpath "$.length()" == 3
jsonpath "$[0].name" == "Laptop"

# 2. Obtener el articulo con id 2
GET http://localhost:8080/articulos/2
HTTP 200
[Asserts]
jsonpath "$.name" == "Smartphone"

# 3. Eliminar el articulo con id 2
DELETE http://localhost:8080/articulos/2
HTTP 204

# 4. Crear un nuevo articulo
POST http://localhost:8080/articulos
Content-Type: application/json
{
"name": "GPS Watch",
"price": 199.05
}
HTTP 201
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.name" == "GPS Watch"

