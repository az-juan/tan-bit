-- name: CreateUsuario :one
INSERT INTO usuario (id_usuario, apellido, nombre, email)
VALUES ($1, $2, $3, $4)
RETURNING id_usuario, apellido, nombre, email;

-- name: DeleteUsuario :exec
DELETE FROM usuario
WHERE id_usuario = $1;

-- name: ListUsuarios :many
SELECT apellido, nombre
FROM usuario
ORDER BY apellido;