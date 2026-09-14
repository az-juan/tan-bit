-- name: GetProductByID :one
-- Obtiene un producto por su clave primaria
SELECT id_producto, titulo, descripcion, condicion, precio, stock, id_vendedor
FROM producto
WHERE id_producto = $1;

-- name: ListProducts :many
-- Lista productos con paginación usando LIMIT y OFFSET
SELECT id_producto, titulo, descripcion, condicion, precio, stock, id_vendedor
FROM producto
ORDER BY id_producto
LIMIT $1 OFFSET $2;

-- name: ListProductsBySeller :many
-- Filtra todos los productos pertenecientes a un vendedor específico
SELECT id_producto, titulo, descripcion, condicion, precio, stock, id_vendedor
FROM producto
WHERE id_vendedor = $1
ORDER BY id_producto;

-- name: SearchProductsByTitle :many
-- Búsqueda de productos por coincidencia parcial en el título (case-insensitive en PostgreSQL)
SELECT id_producto, titulo, descripcion, condicion, precio, stock, id_vendedor
FROM producto
WHERE titulo ILIKE '%' || $1 || '%'
ORDER BY id_producto;

-- name: GetProductWithSellerDetails :one
-- Consulta con JOIN para traer información combinada del producto y de su vendedor/usuario
SELECT
    p.id_producto,
    p.titulo,
    p.descripcion,
    p.condicion,
    p.precio,
    p.stock,
    p.id_vendedor,
    v.telefono AS vendedor_telefono,
    u.nombre AS vendedor_nombre,
    u.apellido AS vendedor_apellido,
    u.email AS vendedor_email
FROM producto p
JOIN vendedor v ON p.id_vendedor = v.id_vendedor
JOIN usuario u ON v.id_vendedor = u.id_usuario
WHERE p.id_producto = $1;

-- name: CreateProduct :one
-- Inserta un nuevo producto y devuelve el registro completo con su id generado
INSERT INTO producto (
    titulo,
    descripcion,
    condicion,
    precio,
    stock,
    id_vendedor
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id_producto, titulo, descripcion, condicion, precio, stock, id_vendedor;

-- name: DeleteProduct :exec
-- Elimina un producto por su clave primaria
DELETE FROM producto
WHERE id_producto = $1;

-- name: DeleteProductByIDAndSeller :execrows
-- Elimina un producto asegurando que pertenezca al vendedor indicado
-- Retorna la cantidad de filas afectadas para comprobar si se eliminó
DELETE FROM producto
WHERE id_producto = $1 AND id_vendedor = $2;
