-- Create "carrito" table
CREATE TABLE "public"."carrito" (
  "id_usuario" integer NOT NULL,
  "id_articulo" integer NOT NULL,
  "cantidad" integer NOT NULL DEFAULT 1,
  "fecha_agregado" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id_usuario", "id_articulo"),
  CONSTRAINT "fk_articulo" FOREIGN KEY ("id_articulo") REFERENCES "public"."articulo" ("id_articulo") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_usuario" FOREIGN KEY ("id_usuario") REFERENCES "public"."usuario" ("id_usuario") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "favorito" table
CREATE TABLE "public"."favorito" (
  "id_usuario" integer NOT NULL,
  "id_articulo" integer NOT NULL,
  "fecha" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id_usuario", "id_articulo"),
  CONSTRAINT "fk_articulo" FOREIGN KEY ("id_articulo") REFERENCES "public"."articulo" ("id_articulo") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_usuario" FOREIGN KEY ("id_usuario") REFERENCES "public"."usuario" ("id_usuario") ON UPDATE NO ACTION ON DELETE NO ACTION
);
