-- name: CreateVendedor :one
INSERT INTO vendedor (id_vendedor, telefono)
VALUES ($1, $2)
RETURNING id_vendedor, telefono;

-- name: DeleteVendedor :exec
DELETE FROM vendedor
WHERE id_vendedor = $1;