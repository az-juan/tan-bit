-- name: GetArticuloID :one
SELECT id, nombre,precio,descripcion,condicion,categoria,stock,contacto,fecha_publicacion
FROM articulo
WHERE id = $1;

-- name: ListArticulos :many
SELECT id, nombre, precio, descripcion, condicion, categoria, stock, contacto
FROM articulo
ORDER BY nombre;

-- name: CreateArticulo :one
INSERT INTO articulo (nombre, precio, descripcion, condicion, categoria, stock, contacto)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, nombre, precio, descripcion, condicion, categoria, stock, contacto, fecha_publicacion;

-- name: UpdateArticulo :exec
UPDATE articulo
SET nombre = $2, precio = $3, descripcion = $4, condicion = $5, categoria = $6, stock = $7, contacto = $8
WHERE id = $1;

-- name: DeleteArticulo :exec
DELETE FROM articulo
WHERE id = $1;
