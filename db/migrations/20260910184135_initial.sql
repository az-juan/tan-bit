-- Create enum type "condicion"
CREATE TYPE "condicion" AS ENUM ('nuevo', 'muy bueno', 'decente', 'a reparar');
-- Create enum type "categoria"
CREATE TYPE "categoria" AS ENUM ('notebook', 'pantalla', 'electronica', 'periferico', 'accesorio', 'audio', 'video', 'gaming');
-- Create "articulo" table
CREATE TABLE "articulo" (
  "id" serial NOT NULL,
  "nombre" character varying(255) NOT NULL,
  "precio" numeric(10,2) NOT NULL,
  "descripcion" text NOT NULL,
  "condicion" "condicion" NOT NULL,
  "ruta_imagen" text NULL,
  "categoria" "categoria" NOT NULL,
  "stock" integer NOT NULL,
  "contacto" character varying(255) NOT NULL,
  "fecha_publicacion" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "articulo_precio_check" CHECK (precio >= (0)::numeric),
  CONSTRAINT "articulo_stock_check" CHECK (stock >= 0)
);
