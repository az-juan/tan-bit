-- name: GetUserByID :one
-- Obtiene un usuario por su clave primaria
SELECT id_usuario, apellido, nombre, email
FROM usuario
WHERE id_usuario = $1;

-- name: ListUsers :many
-- Lista todos los usuarios
SELECT id_usuario, apellido, nombre, email
FROM usuario
ORDER BY id_usuario;

-- name: CreateUser :one
-- Inserta un nuevo usuario y retorna sus datos con el ID asignado
INSERT INTO usuario (
    nombre,
    apellido,
    email
) VALUES (
    $1, $2, $3
)
RETURNING id_usuario, apellido, nombre, email;

-- name: UpdateUser :exec
-- Actualiza los datos de un usuario existente
UPDATE usuario
SET
    nombre = $2,
    apellido = $3,
    email = $4
WHERE id_usuario = $1;

-- name: DeleteUser :exec
-- Elimina un usuario por su clave primaria
DELETE FROM usuario
WHERE id_usuario = $1;