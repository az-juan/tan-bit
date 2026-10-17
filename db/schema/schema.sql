CREATE TYPE condicion AS ENUM ('nuevo', 'muy bueno', 'usado', 'reacondicionado', 'a reparar');
CREATE TYPE categoria AS ENUM ('notebook','pantalla','electronica','periferico','accesorio','audio','video','gaming');

CREATE TABLE usuario (
    id_usuario  SERIAL      CONSTRAINT pk_usuario PRIMARY KEY,
    apellido    VARCHAR(50) NOT NULL,
    nombre      VARCHAR(50) NOT NULL,
    email       VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE vendedor (
    id_vendedor     INTEGER     CONSTRAINT pk_vendedor PRIMARY KEY CONSTRAINT fk_vendedor_usuario REFERENCES usuario(id_usuario),
    telefono        VARCHAR(20) NOT NULL
);

CREATE TABLE articulo (
    id_articulo     SERIAL       CONSTRAINT pk_producto PRIMARY KEY,
    titulo          VARCHAR(50)  NOT NULL,
    descripcion     VARCHAR(250) NOT NULL,
    condicion       condicion    NOT NULL,
    categoria       categoria    NOT NULL,
    ruta_imagen     text         NULL,
    precio          NUMERIC(10, 2) NOT NULL CHECK (precio >= 0),
    stock           INTEGER      NOT NULL CHECK (stock >= 0),
    fecha_publicacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    id_vendedor     INTEGER      NOT NULL CONSTRAINT fk_producto_vendedor REFERENCES vendedor(id_vendedor)
);

CREATE TABLE compra (
    id_compra       SERIAL          CONSTRAINT pk_compra PRIMARY KEY,
    id_comprador    INTEGER         NOT NULL CONSTRAINT fk_compra_usuario REFERENCES usuario(id_usuario),
    fecha           TIMESTAMP       NOT NULL,
    total           NUMERIC(10, 2)  NOT NULL
);

CREATE TABLE detalle_compra (
    id_compra       INTEGER     NOT NULL CONSTRAINT fk_detalle_compra REFERENCES compra(id_compra),
    id_articulo     INTEGER     NOT NULL CONSTRAINT fk_detalle_producto REFERENCES articulo(id_articulo),
    cantidad        INTEGER     NOT NULL,
    precio_historico_unitario NUMERIC(10, 2) NOT NULL,
    CONSTRAINT pk_detalle_compra PRIMARY KEY (id_compra, id_articulo)
);