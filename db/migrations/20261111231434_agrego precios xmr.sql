-- Modify "articulo" table
ALTER TABLE "public"."articulo" ADD CONSTRAINT "articulo_precio_xmr_check" CHECK (precio_xmr >= (0)::double precision), ALTER COLUMN "precio" TYPE double precision, ALTER COLUMN "costo_envio" TYPE double precision, ADD COLUMN "precio_xmr" double precision NOT NULL;
-- Modify "compra" table
ALTER TABLE "public"."compra" ADD CONSTRAINT "compra_total_check1" CHECK (total >= (0)::double precision), ALTER COLUMN "total" TYPE double precision, ADD COLUMN "total_xmr" double precision NOT NULL;
-- Modify "detalle_compra" table
ALTER TABLE "public"."detalle_compra" ADD CONSTRAINT "detalle_compra_precio_historico_unitario_check" CHECK (precio_historico_unitario >= (1)::double precision), ADD CONSTRAINT "detalle_compra_precio_historico_unitario_xmr_check" CHECK (precio_historico_unitario_xmr >= (0)::double precision), ALTER COLUMN "precio_historico_unitario" TYPE double precision, ADD COLUMN "precio_historico_unitario_xmr" double precision NOT NULL;
