-- name: GetSellerByID :one
-- Obtiene el vendedor junto con los datos de contacto del usuario
SELECT
    v.id_vendedor,
    v.telefono,
    u.nombre,
    u.apellido,
    u.email
FROM vendedor v
JOIN usuario u ON v.id_vendedor = u.id_usuario
WHERE v.id_vendedor = $1;

-- name: ListSellers :many
-- Lista todos los vendedores junto con su información de usuario con paginación
SELECT
    v.id_vendedor,
    v.telefono,
    u.nombre,
    u.apellido,
    u.email
FROM vendedor v
JOIN usuario u ON v.id_vendedor = u.id_usuario
ORDER BY v.id_vendedor
LIMIT $1 OFFSET $2;

-- name: CreateSeller :one
-- Convierte a un usuario existente en vendedor registrando su teléfono
INSERT INTO vendedor (
    id_vendedor,
    telefono
) VALUES (
    $1, $2
)
RETURNING id_vendedor, telefono;

-- name: UpdateSeller :one
-- Actualiza el teléfono de un vendedor
UPDATE vendedor
SET telefono = $2
WHERE id_vendedor = $1
RETURNING id_vendedor, telefono;
