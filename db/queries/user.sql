-- name: GetUserByID :one
-- Obtiene un usuario por su clave primaria
SELECT id_usuario, apellido, nombre, email
FROM usuario
WHERE id_usuario = $1;

-- name: GetUserByEmail :one
-- Obtiene un usuario a partir de su correo electrónico
SELECT id_usuario, apellido, nombre, email
FROM usuario
WHERE email = $1;

-- name: ListUsers :many
-- Lista usuarios con paginación
SELECT id_usuario, apellido, nombre, email
FROM usuario
ORDER BY id_usuario
LIMIT $1 OFFSET $2;

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

-- name: UpdateUser :one
-- Actualiza los datos de un usuario existente
UPDATE usuario
SET
    nombre = $2,
    apellido = $3,
    email = $4
WHERE id_usuario = $1
RETURNING id_usuario, apellido, nombre, email;
