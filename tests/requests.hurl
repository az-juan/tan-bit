# =============================================================
# Smoke tests HTTP + Seed de catálogo de tan-bit
# --file-root para que hurl pueda ver las imagenes de static/assets.
# Uso: hurl --test --file-root . --variable host=http://localhost:8080 tests/requests.hurl
# =============================================================

# ---------------- HEALTH ----------------
GET {{host}}/health
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
body contains "EL HARDWARE TAMBIÉN TIENE HISTORIA."

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

GET {{host}}/vender
HTTP 303
[Asserts]
header "Location" == "/login"

# ---------------- REGISTRO DE VENDEDOR (Lionel Messi) ----------------
# 1) Sin aceptar terminos -> 303 redirect a /registrarse (no toca la DB)
POST {{host}}/registrarse
[Form]
role: vendedor
nombre: Lionel
apellido: Messi
email: leo@messi.com
password: campeon2022
password2: campeon2022
HTTP 303
[Asserts]
header "Location" == "/registrarse"

# 2) Rol vendedor sin telefono/cbu/alias -> 303 redirect a /registrarse (no toca la DB)
POST {{host}}/registrarse
[Form]
role: vendedor
nombre: Lionel
apellido: Messi
email: leo@messi.com
password: campeon2022
password2: campeon2022
terminos: on
HTTP 303
[Asserts]
header "Location" == "/registrarse"

# 3) Registro de vendedor OK -> 303 a /mi_cuenta + cookie user_id
POST {{host}}/registrarse
[Form]
role: vendedor
nombre: Lionel
apellido: Messi
email: leo@messi.com
password: campeon2022
password2: campeon2022
telefono: 2494101010
cbu: 0000010101010
alias: leo.campeon.xmr
direccion_xmr: 44Affq5kSiGBoZ4NMDwYtN18obc8AemS33DBLWs3H7otXft3XjrpDtQGv7SqSsaBYBb98uNbr2VBBEt7f2wfn3RVGQBEP3A
terminos: on
HTTP 303
[Captures]
vendedor_id: cookie "user_id"
[Asserts]
header "Location" == "/mi_cuenta"
cookie "user_id" exists

# 4) Email duplicado -> violacion UNIQUE manejada, encola notificacion y redirige a /registrarse
POST {{host}}/registrarse
[Form]
role: vendedor
nombre: Otro
apellido: Usuario
email: leo@messi.com
password: campeon2022
password2: campeon2022
telefono: 2494999999
cbu: 9999999999999
alias: otro.alias
terminos: on
HTTP 303
[Asserts]
header "Location" == "/registrarse"

# Verifica que la sesion del vendedor muestra su perfil, datos de vendedor y la notificacion encolada
GET {{host}}/mi_cuenta
HTTP 200
[Asserts]
body contains "Lionel Messi"
body contains "leo@messi.com"
body contains "■ VENDEDOR"
body contains "leo.campeon.xmr"
body contains "El email ya está registrado."

# ---------------- LOGIN ----------------
POST {{host}}/login
[Form]
email: leo@messi.com
password: incorrecta
HTTP 303
[Asserts]
header "Location" == "/login"

POST {{host}}/login
[Form]
email: leo@messi.com
password: campeon2022
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"
cookie "user_id" exists

# ---------------- PUBLICACION DE ARTICULOS (/vender) ----------------
# Acceso al formulario de venta con sesion de vendedor
GET {{host}}/vender
HTTP 200
[Asserts]
body contains "DALE CIRCULACIÓN A TU HARDWARE."

# Validacion de campos obligatorios faltantes -> 303 redirect a /vender
POST {{host}}/vender
[MultipartFormData]
titulo:
descripcion:
HTTP 303
[Asserts]
header "Location" == "/vender"

# 1. Gabinete (gabinete.png)
POST {{host}}/vender
[MultipartFormData]
titulo: Gabinete Artesanal de Nogal ATX
categoria: gaming
condicion: nuevo
descripcion: Gabinete custom en madera de nogal macizo y cobre. Lo vendo porque necesito liberar la repisa del escritorio, si no no me entran los 8 Balones de Oro.
ubicacion: Tandil, Buenos Aires
precio: 210000
precio_xmr: 0.95
stock: 1
tipo_entrega: ambos
costo_envio: 4500
zona_entrega: Tandil y alrededores
fotos: file,static/assets/gabinete.png; image/png
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"

# 2. Teclado (teclado.png)
POST {{host}}/vender
[MultipartFormData]
titulo: Teclado Mecanico Custom 75%
categoria: periferico
condicion: muy bueno
descripcion: Chasis mecanizado con switches tactiles lubricados a mano y keycaps PBT retro. Ideal para programar o jugar sin fatiga.
ubicacion: Tandil, Buenos Aires
precio: 115000
precio_xmr: 0.52
stock: 2
tipo_entrega: envio a domicilio
costo_envio: 2500
zona_entrega: Todo el pais
fotos: file,static/assets/teclado.png; image/png
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"

# 3. Memoria RAM (ram.png)
POST {{host}}/vender
[MultipartFormData]
titulo: Kit Memoria RAM DDR4 32GB 3600MHz
categoria: electronica
condicion: muy bueno
descripcion: Dos modulos de 16GB con disipador de aluminio anodizado. Testeados con MemTest86 sin errores, latencia CL16.
ubicacion: Tandil, Buenos Aires
precio: 88000
precio_xmr: 0.40
stock: 3
tipo_entrega: ambos
costo_envio: 2000
zona_entrega: Tandil centro
fotos: file,static/assets/ram.png; image/png
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"

# 4. Motherboard (mother.png)
POST {{host}}/vender
[MultipartFormData]
titulo: Motherboard B550 AM4 VRM Disipado
categoria: electronica
condicion: reacondicionado
descripcion: Placa madre socket AM4 con doble slot M.2 NVMe, PCIe 4.0 y capacitores japoneses. BIOS actualizada lista para usar.
ubicacion: Tandil, Buenos Aires
precio: 145000
precio_xmr: 0.65
stock: 1
tipo_entrega: retiro en persona
costo_envio: 0
zona_entrega: Tandil centro
fotos: file,static/assets/mother.png; image/png
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"

# 5. Cooler (cooler.png)
POST {{host}}/vender
[MultipartFormData]
titulo: Cooler CPU Doble Torre 6 Heatpipes
categoria: accesorio
condicion: nuevo
descripcion: Disipador por aire con base y heatpipes de cobre puro. Silencioso a maxima carga, incluye anclajes AM4, AM5 y LGA1700.
ubicacion: Tandil, Buenos Aires
precio: 62000
precio_xmr: 0.28
stock: 4
tipo_entrega: ambos
costo_envio: 2200
zona_entrega: Tandil y zona serrana
fotos: file,static/assets/cooler.png; image/png
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"

# 6. Lampara (lampara.png)
POST {{host}}/vender
[MultipartFormData]
titulo: Lampara de Escritorio Industrial Vintage
categoria: analogico
condicion: nuevo
descripcion: Armada artesanalmente con tuberias recuperadas y filamento calido regulable. Luz suave para jornadas nocturnas de taller.
ubicacion: Tandil, Buenos Aires
precio: 54000
precio_xmr: 0.24
stock: 2
tipo_entrega: retiro en persona
costo_envio: 0
zona_entrega: Tandil centro
fotos: file,static/assets/lampara.png; image/png
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"

# 7. Pad (pad.png)
POST {{host}}/vender
[MultipartFormData]
titulo: Deskpad XL Texturado Costura Reforzada
categoria: periferico
condicion: nuevo
descripcion: Mousepad extendido de 90x40cm con base de caucho antideslizante y superficie microtexturada de control preciso.
ubicacion: Tandil, Buenos Aires
precio: 29000
precio_xmr: 0.13
stock: 5
tipo_entrega: envio a domicilio
costo_envio: 1800
zona_entrega: Todo el pais
fotos: file,static/assets/pad.png; image/png
HTTP 303
[Asserts]
header "Location" == "/mi_cuenta"

# ---------------- VERIFICACION DE CATALOGO E IMAGENES ----------------
# El marketplace muestra los 7 articulos cargados y sus rutas en /uploads/
GET {{host}}/
HTTP 200
[Asserts]
body contains "Gabinete Artesanal de Nogal ATX"
body contains "no me entran los 8 Balones de Oro"
body contains "Teclado Mecanico Custom 75%"
body contains "Kit Memoria RAM DDR4 32GB 3600MHz"
body contains "Motherboard B550 AM4 VRM Disipado"
body contains "Cooler CPU Doble Torre 6 Heatpipes"
body contains "Lampara de Escritorio Industrial Vintage"
body contains "Deskpad XL Texturado Costura Reforzada"

# El filtro por categoria funciona y filtra los perifericos
GET {{host}}/?categoria=periferico
HTTP 200
[Asserts]
body contains "Teclado Mecanico Custom 75%"
body contains "Deskpad XL Texturado Costura Reforzada"
body not contains "Gabinete Artesanal de Nogal ATX"

# El archivo binario subido al volumen /app/uploads se sirve correctamente
GET {{host}}/uploads/{{vendedor_id}}_gabinete_artesanal_de_nogal_atx_1.png
HTTP 200
[Asserts]
header "Content-Type" contains "image/png"

# En /mi_cuenta del vendedor figuran sus piezas publicadas
GET {{host}}/mi_cuenta
HTTP 200
[Asserts]
body contains "Hay piezas publicadas."
body contains "Gabinete Artesanal de Nogal ATX"


# =============================================================
# TESTS ANTERIORES (API REST JSON, comentada en /handlers)
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
