-- Create enum type "entrega"
CREATE TYPE "public"."entrega" AS ENUM ('envio a domicilio', 'retiro en persona', 'ambos');
-- Modify "articulo" table
ALTER TABLE "public"."articulo" ADD CONSTRAINT "articulo_costo_envio_check" CHECK (costo_envio >= (0)::numeric), ALTER COLUMN "tipo_entrega" TYPE "public"."entrega" USING "tipo_entrega"::"public"."entrega", ALTER COLUMN "tipo_entrega" SET NOT NULL;
-- Modify "vendedor" table
ALTER TABLE "public"."vendedor" ADD COLUMN "direccion_xmr" text NULL, ADD COLUMN "cbu" text NOT NULL, ADD COLUMN "alias" text NOT NULL, ADD CONSTRAINT "vendedor_alias_key" UNIQUE ("alias"), ADD CONSTRAINT "vendedor_cbu_key" UNIQUE ("cbu"), ADD CONSTRAINT "vendedor_direccion_xmr_key" UNIQUE ("direccion_xmr"), ADD CONSTRAINT "vendedor_telefono_key" UNIQUE ("telefono");
