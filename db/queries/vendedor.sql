-- name: CreateVendedor :one
INSERT INTO vendedor (id_vendedor, telefono)
VALUES ($1, $2)
RETURNING id_vendedor, telefono;

-- name: GetVendedores :many
SELECT id_vendedor
FROM vendedor
ORDER BY id_vendedor;

-- name: DeleteVendedor :exec
DELETE FROM vendedor
WHERE id_vendedor = $1;