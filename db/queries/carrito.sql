-- name: AddToCarrito :exec
INSERT INTO carrito (id_usuario, id_articulo, cantidad)
VALUES ($1, $2, $3)
ON CONFLICT (id_usuario, id_articulo) DO UPDATE SET cantidad = carrito.cantidad + EXCLUDED.cantidad;

-- name: GetCarritoByUsuario :many
SELECT a.id_articulo, a.titulo, a.descripcion, a.condicion, a.categoria, a.precio, a.stock, a.fecha_publicacion, a.ruta_imagen, c.cantidad
FROM carrito c JOIN articulo a ON a.id_articulo = c.id_articulo
WHERE c.id_usuario = $1
ORDER BY c.fecha_agregado;

-- name: RemoveFromCarrito :exec
DELETE FROM carrito WHERE id_usuario = $1 AND id_articulo = $2;
