-- name: GetArticuloByID :one
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, precio_xmr, stock, fecha_publicacion, ruta_imagen, ubicacion, tipo_entrega, costo_envio, zona_entrega
FROM articulo
WHERE id_articulo = $1;

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

-- name: UpdateArticulo :execrows
UPDATE articulo
SET titulo = $2, descripcion = $3, condicion = $4, categoria = $5, precio = $6, precio_xmr = $7, stock = $8, ubicacion = $9, tipo_entrega = $10, costo_envio = $11, zona_entrega = $12
WHERE id_articulo = $1;

-- name: DeleteArticulo :execrows
DELETE FROM articulo
WHERE id_articulo = $1;

-- name: GetArticulosByVendedor :many
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, precio_xmr, stock, fecha_publicacion, ruta_imagen, ubicacion, tipo_entrega, costo_envio, zona_entrega, id_vendedor
FROM articulo
WHERE id_vendedor = $1
ORDER BY fecha_publicacion DESC;
