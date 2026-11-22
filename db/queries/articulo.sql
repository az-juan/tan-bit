-- name: GetArticuloByID :one
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, precio_xmr, stock, fecha_publicacion, ruta_imagen, ubicacion, tipo_entrega, costo_envio, zona_entrega
FROM articulo
WHERE id_articulo = $1;

-- name: GetArticuloCompletoByID :one
SELECT a.id_articulo, a.titulo, a.descripcion, a.condicion, a.categoria, a.precio, a.precio_xmr, a.stock, a.fecha_publicacion, a.ruta_imagen, a.ubicacion, a.tipo_entrega, a.costo_envio, a.zona_entrega, v.alias, v.cbu, v.direccion_xmr
FROM articulo a
JOIN vendedor v ON v.id_vendedor = a.id_vendedor
WHERE a.id_articulo = $1;

-- name: GetArticuloCompletoByVendedor :one
SELECT a.id_articulo, a.titulo, a.descripcion, a.condicion, a.categoria, a.precio, a.precio_xmr, a.stock, a.fecha_publicacion, a.ruta_imagen, a.ubicacion, a.tipo_entrega, a.costo_envio, a.zona_entrega, v.alias, v.cbu, v.direccion_xmr
FROM articulo a
JOIN vendedor v ON v.id_vendedor = a.id_vendedor
WHERE a.id_articulo = $1 AND a.id_vendedor = $2;

-- name: DecrementarStock :one
UPDATE articulo
SET stock = stock - $2
WHERE id_articulo = $1 AND stock >= $2
RETURNING id_articulo, stock;

-- name: GetArticulos :many
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, precio_xmr, stock, fecha_publicacion, ruta_imagen, ubicacion, tipo_entrega, costo_envio, zona_entrega
FROM articulo
ORDER BY id_articulo;

-- name: CreateArticulo :one
INSERT INTO articulo (titulo, descripcion, condicion, categoria, precio, precio_xmr, stock, id_vendedor, ubicacion, tipo_entrega, costo_envio, zona_entrega, ruta_imagen)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING id_articulo, titulo, descripcion, condicion, categoria, precio, precio_xmr, stock, id_vendedor, ubicacion, tipo_entrega, costo_envio, zona_entrega, ruta_imagen;

-- name: GetArticulosByCategoria :many
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, precio_xmr, stock, fecha_publicacion, ruta_imagen, ubicacion, tipo_entrega, costo_envio, zona_entrega
FROM articulo
WHERE categoria = $1
ORDER BY id_articulo;

-- name: SearchArticulos :many
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, precio_xmr, stock, fecha_publicacion, ruta_imagen, ubicacion, tipo_entrega, costo_envio, zona_entrega
FROM articulo
WHERE titulo ILIKE '%' || $1 || '%' OR descripcion ILIKE '%' || $1 || '%'
ORDER BY id_articulo;

-- name: UpdateArticulo :execrows
UPDATE articulo
SET titulo = $2, descripcion = $3, condicion = $4, categoria = $5, precio = $6, precio_xmr = $7, stock = $8, ubicacion = $9, tipo_entrega = $10, costo_envio = $11, zona_entrega = $12
WHERE id_articulo = $1;

-- name: UpdateArticuloDeVendedor :execrows
UPDATE articulo
SET titulo = $3, descripcion = $4, condicion = $5, categoria = $6, precio = $7, precio_xmr = $8, stock = $9, ubicacion = $10, tipo_entrega = $11, costo_envio = $12, zona_entrega = $13
WHERE id_articulo = $1 AND id_vendedor = $2;

-- name: DeleteArticulo :execrows
DELETE FROM articulo
WHERE id_articulo = $1;

-- name: DeleteArticuloDeVendedor :execrows
DELETE FROM articulo
WHERE id_articulo = $1 AND id_vendedor = $2;

-- name: GetArticulosByVendedor :many
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, precio_xmr, stock, fecha_publicacion, ruta_imagen, ubicacion, tipo_entrega, costo_envio, zona_entrega, id_vendedor
FROM articulo
WHERE id_vendedor = $1
ORDER BY fecha_publicacion DESC;
