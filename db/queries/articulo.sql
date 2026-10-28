-- name: GetArticuloByID :one
SELECT titulo, descripcion, condicion, categoria, precio, stock, fecha_publicacion
FROM articulo
WHERE id_articulo = $1;

-- name: GetArticulos :many
SELECT titulo, descripcion, condicion, categoria, precio, stock, fecha_publicacion
FROM articulo
ORDER BY id_articulo;

-- name: CreateArticulo :one
INSERT INTO articulo (titulo, descripcion, condicion, categoria, precio, stock, id_vendedor)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING titulo, descripcion, condicion, categoria, precio, stock, id_vendedor;

-- name: UpdateArticulo :exec
UPDATE articulo
SET titulo = $2, descripcion = $3, condicion = $4, precio = $5, stock = $6
WHERE id_articulo = $1;

-- name: DeleteArticulo :exec
DELETE FROM articulo
WHERE id_articulo = $1;
