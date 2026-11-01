-- name: GetArticuloByID :one
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, stock, fecha_publicacion
FROM articulo
WHERE id_articulo = $1;

-- name: GetArticulos :many
SELECT id_articulo, titulo, descripcion, condicion, categoria, precio, stock, fecha_publicacion
FROM articulo
ORDER BY id_articulo;

-- name: CreateArticulo :one
INSERT INTO articulo (titulo, descripcion, condicion, categoria, precio, stock, id_vendedor)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id_articulo, titulo, descripcion, condicion, categoria, precio, stock, id_vendedor;

-- name: UpdateArticulo :execrows
UPDATE articulo
SET titulo = $2, descripcion = $3, condicion = $4, categoria = $5, precio = $6, stock = $7
WHERE id_articulo = $1;

-- name: DeleteArticulo :execrows
DELETE FROM articulo
WHERE id_articulo = $1;