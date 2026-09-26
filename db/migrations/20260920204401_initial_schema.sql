-- Create enum type "condicion_enum"
CREATE TYPE "condicion_enum" AS ENUM ('nuevo', 'usado', 'reacondicionado', 'roto');
-- Create "usuario" table
CREATE TABLE "usuario" (
  "id_usuario" serial NOT NULL,
  "apellido" character varying(50) NOT NULL,
  "nombre" character varying(50) NOT NULL,
  "email" character varying(50) NOT NULL,
  CONSTRAINT "pk_usuario" PRIMARY KEY ("id_usuario"),
  CONSTRAINT "usuario_email_key" UNIQUE ("email")
);
-- Create "compra" table
CREATE TABLE "compra" (
  "id_compra" serial NOT NULL,
  "id_comprador" integer NOT NULL,
  "fecha" timestamp NOT NULL,
  "total" numeric(10,2) NOT NULL,
  CONSTRAINT "pk_compra" PRIMARY KEY ("id_compra"),
  CONSTRAINT "fk_compra_usuario" FOREIGN KEY ("id_comprador") REFERENCES "usuario" ("id_usuario") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "vendedor" table
CREATE TABLE "vendedor" (
  "id_vendedor" integer NOT NULL,
  "telefono" character varying(20) NOT NULL,
  CONSTRAINT "pk_vendedor" PRIMARY KEY ("id_vendedor"),
  CONSTRAINT "fk_vendedor_usuario" FOREIGN KEY ("id_vendedor") REFERENCES "usuario" ("id_usuario") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "producto" table
CREATE TABLE "producto" (
  "id_producto" serial NOT NULL,
  "titulo" character varying(50) NOT NULL,
  "descripcion" character varying(250) NOT NULL,
  "condicion" "condicion_enum" NOT NULL,
  "precio" numeric(10,2) NOT NULL,
  "stock" integer NOT NULL,
  "id_vendedor" integer NOT NULL,
  CONSTRAINT "pk_producto" PRIMARY KEY ("id_producto"),
  CONSTRAINT "fk_producto_vendedor" FOREIGN KEY ("id_vendedor") REFERENCES "vendedor" ("id_vendedor") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "detalle_compra" table
CREATE TABLE "detalle_compra" (
  "id_compra" integer NOT NULL,
  "id_producto" integer NOT NULL,
  "cantidad" integer NOT NULL,
  "precio_historico_unitario" numeric(10,2) NOT NULL,
  CONSTRAINT "pk_detalle_compra" PRIMARY KEY ("id_compra", "id_producto"),
  CONSTRAINT "fk_detalle_compra" FOREIGN KEY ("id_compra") REFERENCES "compra" ("id_compra") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_detalle_producto" FOREIGN KEY ("id_producto") REFERENCES "producto" ("id_producto") ON UPDATE NO ACTION ON DELETE NO ACTION
);
