-- name: GetArticuloByID :one
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, stock, fecha_publicacion, ruta_imagen, ubicacion
FROM articulo
WHERE id_articulo = $1;

-- name: GetArticulos :many
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, stock, fecha_publicacion, ruta_imagen, ubicacion
FROM articulo
ORDER BY id_articulo;

-- name: CreateArticulo :one
INSERT INTO articulo (titulo, descripcion, condicion, categoria, precio, stock, id_vendedor, ubicacion)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id_articulo, titulo, descripcion, condicion, categoria, precio, stock, id_vendedor, ubicacion;

-- name: GetArticulosByCategoria :many
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, stock, fecha_publicacion, ruta_imagen, ubicacion
FROM articulo
WHERE categoria = $1
ORDER BY id_articulo;

-- name: UpdateArticulo :execrows
UPDATE articulo
SET titulo = $2, descripcion = $3, condicion = $4, categoria = $5, precio = $6, stock = $7, ubicacion = $8
WHERE id_articulo = $1;

-- name: DeleteArticulo :execrows
DELETE FROM articulo
WHERE id_articulo = $1;
