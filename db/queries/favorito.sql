-- name: AddFavorito :exec
INSERT INTO favorito (id_usuario, id_articulo)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: GetFavoritosByUsuario :many
SELECT a.id_articulo, a.titulo, a.descripcion, a.condicion, a.categoria, a.precio, a.precio_xmr, a.stock, a.fecha_publicacion, a.ruta_imagen
FROM favorito f JOIN articulo a ON a.id_articulo = f.id_articulo
WHERE f.id_usuario = $1
ORDER BY f.fecha;

-- name: RemoveFavorito :exec
DELETE FROM favorito WHERE id_usuario = $1 AND id_articulo = $2;
