-- CREATE DOMAIN u_int AS integer
-- CHECK (VALUE >= 0);

CREATE TYPE condicion AS ENUM ('nuevo','muy bueno','decente','a reparar');
CREATE TYPE categoria AS ENUM ('notebook','pantalla','electronica','periferico','accesorio','audio','video','gaming');
-- categoria creada por los usuarios?

CREATE TABLE articulo (
    id          SERIAL        PRIMARY KEY,
    nombre      varchar(255)  NOT NULL,
    precio      decimal(10,2) NOT NULL CHECK (precio >= 0),
    descripcion text          NOT NULL,
    condicion   condicion     NOT NULL,
    ruta_imagen text              NULL,
    categoria   categoria     NOT NULL,
    stock       integer       NOT NULL CHECK (stock >= 0),
    contacto    varchar(255)  NOT NULL,
    fecha_publicacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
