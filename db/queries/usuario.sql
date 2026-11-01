-- name: CreateUsuario :one
INSERT INTO usuario (apellido, nombre, email)
VALUES ($1, $2, $3)
RETURNING id_usuario, apellido, nombre, email;

-- name: DeleteUsuario :execrows
DELETE FROM usuario
WHERE id_usuario = $1;

-- name: GetUsuarios :many
SELECT id_usuario, apellido, nombre, email
FROM usuario
ORDER BY apellido;

-- name: GetUsuarioByID :one
SELECT id_usuario, apellido, nombre, email
FROM usuario
WHERE id_usuario = $1;

-- name: UpdateUsuario :execrows
UPDATE usuario
SET apellido = $2, nombre = $3, email = $4
WHERE id_usuario = $1;
