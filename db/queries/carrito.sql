-- name: AddToCarrito :exec
INSERT INTO carrito (id_usuario, id_articulo, cantidad)
VALUES ($1, $2, $3)
ON CONFLICT (id_usuario, id_articulo) DO UPDATE SET cantidad = carrito.cantidad + EXCLUDED.cantidad;

-- name: GetCarritoCompra :many
SELECT a.id_articulo, a.titulo, a.descripcion, a.condicion, a.categoria, a.precio, a.precio_xmr, a.stock, a.fecha_publicacion, a.ruta_imagen, a.tipo_entrega, v.alias, v.cbu, v.direccion_xmr, c.cantidad
FROM carrito c
JOIN articulo a ON a.id_articulo = c.id_articulo
JOIN vendedor v ON v.id_vendedor = a.id_vendedor
WHERE c.id_usuario = $1
ORDER BY c.fecha_agregado;

-- name: GetCarritoByUsuario :many
SELECT a.id_articulo, a.titulo, a.descripcion, a.condicion, a.categoria, a.precio, a.precio_xmr, a.stock, a.fecha_publicacion, a.ruta_imagen, c.cantidad
FROM carrito c JOIN articulo a ON a.id_articulo = c.id_articulo
WHERE c.id_usuario = $1
ORDER BY c.fecha_agregado;

-- name: RemoveFromCarrito :exec
DELETE FROM carrito WHERE id_usuario = $1 AND id_articulo = $2;
