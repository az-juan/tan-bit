-- name: CreateCompra :one
INSERT INTO compra (id_comprador, fecha, total, total_xmr)
VALUES ($1, NOW(), $2, $3)
RETURNING id_compra, id_comprador, fecha, total, total_xmr;

-- name: CreateDetalleCompra :exec
INSERT INTO detalle_compra (id_compra, id_articulo, cantidad, precio_historico_unitario, precio_historico_unitario_xmr)
VALUES ($1, $2, $3, $4, $5);

-- name: ClearCarrito :exec
DELETE FROM carrito WHERE id_usuario = $1;

-- name: GetComprasByUsuario :many
SELECT id_compra, id_comprador, fecha, total
FROM compra
WHERE id_comprador = $1
ORDER BY fecha DESC;
