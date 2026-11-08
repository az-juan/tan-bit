-- Modify "articulo" table
ALTER TABLE "public"."articulo" ADD COLUMN "tipo_entrega" text NULL, ADD COLUMN "costo_envio" numeric(10,2) NULL, ADD COLUMN "zona_entrega" text NULL;
