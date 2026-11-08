-- name: CreateUsuario :one
INSERT INTO usuario (apellido, nombre, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id_usuario, apellido, nombre, email, password;

-- name: DeleteUsuario :execrows
DELETE FROM usuario
WHERE id_usuario = $1;

-- name: GetUsuarios :many
SELECT id_usuario, apellido, nombre, email, password
FROM usuario
ORDER BY apellido;

-- name: GetUsuarioByEmail :one
SELECT id_usuario, apellido, nombre, email, password
FROM usuario
WHERE email = $1;

-- name: GetUsuarioByID :one
SELECT id_usuario, apellido, nombre, email, password
FROM usuario
WHERE id_usuario = $1;

-- name: UpdateUsuario :execrows
UPDATE usuario
SET apellido = $2, nombre = $3, email = $4, password = $5
WHERE id_usuario = $1;