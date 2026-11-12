# =============================================================
# Smoke tests HTTP de tan-bit (vistas HTML + flujo de sesion)
# Uso: hurl --test --variable host=http://localhost:8080 tests/hurl_test.hurl
# Nota: no hay endpoint de borrado -> cada corrida deja 1 usuario en la DB.
# =============================================================

# ---------------- HEALTH ----------------
GET {{host}}/health
[Options]
variable: email=hurl_{{newUuid}}@test.io
HTTP 200
[Asserts]
body == "Servidor andando!"

POST {{host}}/health
HTTP 405

# ---------------- PAGINAS PUBLICAS ----------------
GET {{host}}/
HTTP 200
[Asserts]
body contains "TAN-BIT · Marketplace"

GET {{host}}/historia
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "TAN-BIT · Historia"

GET {{host}}/static/css/shared.css
HTTP 200

GET {{host}}/no-existe
HTTP 404
[Asserts]
body contains "es inválida"

# ---------------- SIN SESION -> redirect a /login ----------------
GET {{host}}/mi_cuenta
HTTP 303
[Asserts]
header "Location" == "/login"

GET {{host}}/carrito
HTTP 303
[Asserts]
header "Location" == "/login"

# ---------------- REGISTRO ----------------
# Sin aceptar terminos -> re-renderiza el form con error (no toca la DB)
POST {{host}}/registrarse
[Form]
nombre: Lionel
apellido: Messi
email: {{email}}
password: hurl1234
HTTP 200
[Asserts]
body contains "Debes aceptar los términos"

# Registro OK -> 303 + cookie de sesion
POST {{host}}/registrarse
[Form]
nombre: Lionel
apellido: Messi
email: {{email}}
password: hurl1234
terminos: on
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"
cookie "user_id" exists

# Email duplicado -> violacion UNIQUE manejada
POST {{host}}/registrarse
[Form]
nombre: Otro
apellido: Usuario
email: {{email}}
password: hurl1234
terminos: on
HTTP 200
[Asserts]
body contains "El email ya está registrado"

# Hurl reenvia la cookie sola -> ve su propia cuenta
GET {{host}}/mi_cuenta
HTTP 200
[Asserts]
body contains "{{email}}"

# ---------------- LOGIN ----------------
POST {{host}}/login
[Form]
email: {{email}}
password: incorrecta
HTTP 401
[Asserts]
body contains "Credenciales inválidas"

POST {{host}}/login
[Form]
email: {{email}}
password: hurl1234
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"
cookie "user_id" exists


# =============================================================
# TESTS ANTERIORES (API REST JSON, hoy comentada en /handlers)
# =============================================================
# # =============================================================
# # Tests de contrato HTTP de tan-bit (status + body)
# # Uso: hurl --test --variable host=http://localhost:8080 tests/hurl_test.hurl
# # =============================================================
#
# # ---------------- HEALTH ----------------
# GET {{host}}/health
# HTTP 200
# [Asserts]
# body == "Servidor andando!"
#
#
# # ================ USUARIOS ================
#
# # Create
# POST {{host}}/usuarios
# {
#   "apellido": "Messi",
#   "nombre": "Leonel",
#   "email": "leo@hurl.test"
# }
# HTTP 201
# [Captures]
# user_id: jsonpath "$.id_usuario"
# [Asserts]
# header "Content-Type" == "application/json"
# jsonpath "$.id_usuario" isInteger
# jsonpath "$.apellido" == "Messi"
# jsonpath "$.nombre" == "Leonel"
# jsonpath "$.email" == "leo@hurl.test"
#
# # Create con email duplicado -> 409
# POST {{host}}/usuarios
# {
#   "apellido": "Otro",
#   "nombre": "Usuario",
#   "email": "leo@hurl.test"
# }
# HTTP 409
# [Asserts]
# body contains "Email ya existente"
#
# # Create con JSON invalido -> 400
# POST {{host}}/usuarios
# HTTP 400
# [Asserts]
# body contains "JSON inválido"
#
# # GetByID
# GET {{host}}/usuarios/{{user_id}}
# HTTP 200
# [Asserts]
# jsonpath "$.id_usuario" == {{user_id}}
# jsonpath "$.email" == "leo@hurl.test"
#
# # GetByID inexistente -> 404
# GET {{host}}/usuarios/999999
# HTTP 404
#
# # GetByID con id no numerico -> 400
# GET {{host}}/usuarios/abc
# HTTP 400
# [Asserts]
# body contains "ID de usuario invalido"
#
# # Get (listado)
# GET {{host}}/usuarios
# HTTP 200
# [Asserts]
# jsonpath "$" isCollection
# jsonpath "$[?(@.id_usuario == {{user_id}})]" exists
#
# # Update
# PUT {{host}}/usuarios/{{user_id}}
# {
#   "apellido": "Messi Cuccittini",
#   "nombre": "Lionel",
#   "email": "lio@hurl.test"
# }
# HTTP 204
#
# # Verificar Update
# GET {{host}}/usuarios/{{user_id}}
# HTTP 200
# [Asserts]
# jsonpath "$.apellido" == "Messi Cuccittini"
# jsonpath "$.nombre" == "Lionel"
# jsonpath "$.email" == "lio@hurl.test"
#
# # Update inexistente -> 404
# PUT {{host}}/usuarios/999999
# {
#   "apellido": "X",
#   "nombre": "Y",
#   "email": "nadie@hurl.test"
# }
# HTTP 404
#
# # Metodo no permitido -> 405
# PATCH {{host}}/usuarios
# HTTP 405
#
#
# # ================ VENDEDORES ================
#
# # Create (el usuario pasa a ser vendedor)
# POST {{host}}/vendedores
# {
#   "id_vendedor": {{user_id}},
#   "telefono": 2494777777
# }
# HTTP 201
# [Asserts]
# jsonpath "$.id_vendedor" == {{user_id}}
# jsonpath "$.telefono" == 2494777777
#
# # Create duplicado (PK) -> 409
# POST {{host}}/vendedores
# {
#   "id_vendedor": {{user_id}},
#   "telefono": 2494000000
# }
# HTTP 409
# [Asserts]
# body contains "El usuario ya es vendedor"
#
# # Create con usuario inexistente (FK) -> 400
# POST {{host}}/vendedores
# {
#   "id_vendedor": 999999,
#   "telefono": 2494000000
# }
# HTTP 400
# [Asserts]
# body contains "No existe tal usuario"
#
# # GetByID
# GET {{host}}/vendedores/{{user_id}}
# HTTP 200
# [Asserts]
# jsonpath "$.id_vendedor" == {{user_id}}
# jsonpath "$.telefono" == 2494777777
#
# # GetByID inexistente -> 404
# GET {{host}}/vendedores/999999
# HTTP 404
#
# # Get (listado)
# GET {{host}}/vendedores
# HTTP 200
# [Asserts]
# jsonpath "$" isCollection
# jsonpath "$[?(@.id_vendedor == {{user_id}})]" exists
#
# # Update
# PUT {{host}}/vendedores/{{user_id}}
# {
#   "telefono": 2494111111
# }
# HTTP 204
#
# # Verificar Update
# GET {{host}}/vendedores/{{user_id}}
# HTTP 200
# [Asserts]
# jsonpath "$.telefono" == 2494111111
#
# # Update inexistente -> 404
# PUT {{host}}/vendedores/999999
# {
#   "telefono": 2494111111
# }
# HTTP 404
#
# # Borrar un usuario que es vendedor -> 409 (FK)
# DELETE {{host}}/usuarios/{{user_id}}
# HTTP 409
#
#
# # ================ ARTICULOS ================
#
# # Create
# POST {{host}}/articulos
# {
#   "titulo": "Thinkpad T14",
#   "descripcion": "notebook empresarial de buena calidad",
#   "condicion": "muy bueno",
#   "categoria": "notebook",
#   "precio": "300000.00",
#   "stock": 1,
#   "id_vendedor": {{user_id}}
# }
# HTTP 201
# [Captures]
# art_id: jsonpath "$.id_articulo"
# [Asserts]
# jsonpath "$.id_articulo" isInteger
# jsonpath "$.titulo" == "Thinkpad T14"
# jsonpath "$.condicion" == "muy bueno"
# jsonpath "$.categoria" == "notebook"
# jsonpath "$.precio" == "300000.00"
# jsonpath "$.stock" == 1
# jsonpath "$.id_vendedor" == {{user_id}}
#
# # Create con vendedor inexistente (FK) -> 400
# POST {{host}}/articulos
# {
#   "titulo": "Huerfano",
#   "descripcion": "sin vendedor",
#   "condicion": "usado",
#   "categoria": "audio",
#   "precio": "10.00",
#   "stock": 1,
#   "id_vendedor": 999999
# }
# HTTP 400
# [Asserts]
# body contains "No existe tal vendedor"
#
# # GetByID
# GET {{host}}/articulos/{{art_id}}
# HTTP 200
# [Asserts]
# jsonpath "$.id_articulo" == {{art_id}}
# jsonpath "$.titulo" == "Thinkpad T14"
# jsonpath "$.fecha_publicacion" isString
#
# # GetByID inexistente -> 404
# GET {{host}}/articulos/999999
# HTTP 404
#
# # Get (listado)
# GET {{host}}/articulos
# HTTP 200
# [Asserts]
# jsonpath "$" isCollection
# jsonpath "$[?(@.id_articulo == {{art_id}})]" exists
#
# # Update
# PUT {{host}}/articulos/{{art_id}}
# {
#   "titulo": "Thinkpad T14 Gen 2",
#   "descripcion": "bateria nueva",
#   "condicion": "usado",
#   "categoria": "notebook",
#   "precio": "250000.00",
#   "stock": 2
# }
# HTTP 204
#
# # Verificar Update
# GET {{host}}/articulos/{{art_id}}
# HTTP 200
# [Asserts]
# jsonpath "$.titulo" == "Thinkpad T14 Gen 2"
# jsonpath "$.condicion" == "usado"
# jsonpath "$.precio" == "250000.00"
# jsonpath "$.stock" == 2
#
# # Update inexistente -> 404
# PUT {{host}}/articulos/999999
# {
#   "titulo": "X",
#   "descripcion": "X",
#   "condicion": "usado",
#   "categoria": "audio",
#   "precio": "1.00",
#   "stock": 1
# }
# HTTP 404
#
# # Borrar vendedor con articulos -> 409 (FK)
# DELETE {{host}}/vendedores/{{user_id}}
# HTTP 409
# [Asserts]
# body contains "El vendedor tiene articulos publicados"
#
#
# # ================ DELETES (orden inverso a las FK) ================
#
# DELETE {{host}}/articulos/{{art_id}}
# HTTP 204
#
# GET {{host}}/articulos/{{art_id}}
# HTTP 404
#
# DELETE {{host}}/articulos/{{art_id}}
# HTTP 404
#
# DELETE {{host}}/vendedores/{{user_id}}
# HTTP 204
#
# GET {{host}}/vendedores/{{user_id}}
# HTTP 404
#
# DELETE {{host}}/usuarios/{{user_id}}
# HTTP 204
#
# GET {{host}}/usuarios/{{user_id}}
# HTTP 404
