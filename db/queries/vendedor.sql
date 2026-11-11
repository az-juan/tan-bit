-- name: CreateVendedor :one
INSERT INTO vendedor (id_vendedor, telefono, direccion_xmr, cbu, alias)
VALUES ($1, $2, $3, $4, $5)
RETURNING id_vendedor, telefono, direccion_xmr, cbu, alias;

-- name: GetVendedores :many
SELECT id_vendedor, telefono, direccion_xmr, cbu, alias
FROM vendedor
ORDER BY id_vendedor;

-- name: DeleteVendedor :execrows
DELETE FROM vendedor
WHERE id_vendedor = $1;

-- name: GetVendedorByID :one
SELECT id_vendedor, telefono, direccion_xmr, cbu, alias
FROM vendedor
WHERE id_vendedor = $1;

-- name: UpdateVendedor :execrows
UPDATE vendedor
SET telefono = $2, direccion_xmr = $3, cbu = $4, alias = $5
WHERE id_vendedor = $1;
