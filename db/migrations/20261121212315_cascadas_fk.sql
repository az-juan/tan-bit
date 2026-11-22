-- Modify "vendedor" table
ALTER TABLE "vendedor" DROP CONSTRAINT "fk_vendedor_usuario", ADD CONSTRAINT "fk_vendedor_usuario" FOREIGN KEY ("id_vendedor") REFERENCES "usuario" ("id_usuario") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "articulo" table
ALTER TABLE "articulo" DROP CONSTRAINT "fk_producto_vendedor", ADD CONSTRAINT "fk_producto_vendedor" FOREIGN KEY ("id_vendedor") REFERENCES "vendedor" ("id_vendedor") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "carrito" table
ALTER TABLE "carrito" DROP CONSTRAINT "fk_articulo", DROP CONSTRAINT "fk_usuario", ADD CONSTRAINT "fk_articulo" FOREIGN KEY ("id_articulo") REFERENCES "articulo" ("id_articulo") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_usuario" FOREIGN KEY ("id_usuario") REFERENCES "usuario" ("id_usuario") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "compra" table
ALTER TABLE "compra" DROP CONSTRAINT "fk_compra_usuario", ADD CONSTRAINT "fk_compra_usuario" FOREIGN KEY ("id_comprador") REFERENCES "usuario" ("id_usuario") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "detalle_compra" table
ALTER TABLE "detalle_compra" DROP CONSTRAINT "fk_detalle_compra", DROP CONSTRAINT "fk_detalle_producto", ADD CONSTRAINT "fk_detalle_compra" FOREIGN KEY ("id_compra") REFERENCES "compra" ("id_compra") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_detalle_producto" FOREIGN KEY ("id_articulo") REFERENCES "articulo" ("id_articulo") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "favorito" table
ALTER TABLE "favorito" DROP CONSTRAINT "fk_articulo", DROP CONSTRAINT "fk_usuario", ADD CONSTRAINT "fk_articulo" FOREIGN KEY ("id_articulo") REFERENCES "articulo" ("id_articulo") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "fk_usuario" FOREIGN KEY ("id_usuario") REFERENCES "usuario" ("id_usuario") ON UPDATE CASCADE ON DELETE CASCADE;
