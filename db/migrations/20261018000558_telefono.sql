-- Create enum type "condicion"
CREATE TYPE "public"."condicion" AS ENUM ('nuevo', 'muy bueno', 'usado', 'reacondicionado', 'a reparar');
-- Create enum type "categoria"
CREATE TYPE "public"."categoria" AS ENUM ('notebook', 'pantalla', 'electronica', 'periferico', 'accesorio', 'audio', 'video', 'gaming');
-- Create "usuario" table
CREATE TABLE "public"."usuario" (
  "id_usuario" serial NOT NULL,
  "apellido" character varying(50) NOT NULL,
  "nombre" character varying(50) NOT NULL,
  "email" character varying(50) NOT NULL,
  CONSTRAINT "pk_usuario" PRIMARY KEY ("id_usuario"),
  CONSTRAINT "usuario_email_key" UNIQUE ("email")
);
-- Create "vendedor" table
CREATE TABLE "public"."vendedor" (
  "id_vendedor" integer NOT NULL,
  "telefono" numeric(10,2) NOT NULL,
  CONSTRAINT "pk_vendedor" PRIMARY KEY ("id_vendedor"),
  CONSTRAINT "fk_vendedor_usuario" FOREIGN KEY ("id_vendedor") REFERENCES "public"."usuario" ("id_usuario") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "articulo" table
CREATE TABLE "public"."articulo" (
  "id_articulo" serial NOT NULL,
  "titulo" character varying(50) NOT NULL,
  "descripcion" character varying(250) NOT NULL,
  "condicion" "public"."condicion" NOT NULL,
  "categoria" "public"."categoria" NOT NULL,
  "ruta_imagen" text NULL,
  "precio" numeric(10,2) NOT NULL,
  "stock" integer NOT NULL,
  "fecha_publicacion" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "id_vendedor" integer NOT NULL,
  CONSTRAINT "pk_producto" PRIMARY KEY ("id_articulo"),
  CONSTRAINT "fk_producto_vendedor" FOREIGN KEY ("id_vendedor") REFERENCES "public"."vendedor" ("id_vendedor") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "articulo_precio_check" CHECK (precio >= (0)::numeric),
  CONSTRAINT "articulo_stock_check" CHECK (stock >= 0)
);
-- Create "compra" table
CREATE TABLE "public"."compra" (
  "id_compra" serial NOT NULL,
  "id_comprador" integer NOT NULL,
  "fecha" timestamp NOT NULL,
  "total" numeric(10,2) NOT NULL,
  CONSTRAINT "pk_compra" PRIMARY KEY ("id_compra"),
  CONSTRAINT "fk_compra_usuario" FOREIGN KEY ("id_comprador") REFERENCES "public"."usuario" ("id_usuario") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "detalle_compra" table
CREATE TABLE "public"."detalle_compra" (
  "id_compra" integer NOT NULL,
  "id_articulo" integer NOT NULL,
  "cantidad" integer NOT NULL,
  "precio_historico_unitario" numeric(10,2) NOT NULL,
  CONSTRAINT "pk_detalle_compra" PRIMARY KEY ("id_compra", "id_articulo"),
  CONSTRAINT "fk_detalle_compra" FOREIGN KEY ("id_compra") REFERENCES "public"."compra" ("id_compra") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_detalle_producto" FOREIGN KEY ("id_articulo") REFERENCES "public"."articulo" ("id_articulo") ON UPDATE NO ACTION ON DELETE NO ACTION
);
