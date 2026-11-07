-- Add value to enum type: "categoria"
ALTER TYPE "public"."categoria" ADD VALUE 'comunicacion';
-- Add value to enum type: "categoria"
ALTER TYPE "public"."categoria" ADD VALUE 'analogico';
-- Add value to enum type: "categoria"
ALTER TYPE "public"."categoria" ADD VALUE 'decoracion';
-- Modify "articulo" table
ALTER TABLE "public"."articulo" ADD COLUMN "ubicacion" text NOT NULL;
-- Modify "carrito" table
ALTER TABLE "public"."carrito" ADD CONSTRAINT "carrito_cantidad_check" CHECK (cantidad >= 1), ALTER COLUMN "cantidad" DROP DEFAULT;
-- Modify "compra" table
ALTER TABLE "public"."compra" ADD CONSTRAINT "compra_total_check" CHECK (total >= (1)::numeric);
-- Modify "detalle_compra" table
ALTER TABLE "public"."detalle_compra" ADD CONSTRAINT "detalle_compra_cantidad_check" CHECK (cantidad >= 1);
