-- name: CreateUsuario :one
INSERT INTO usuario (apellido, nombre, email)
VALUES ($1, $2, $3)
RETURNING apellido, nombre, email;

-- name: DeleteUsuario :exec
DELETE FROM usuario
WHERE id_usuario = $1;

-- name: GetUsuarios :many
SELECT id_usuario, apellido, nombre
FROM usuario
ORDER BY apellido;