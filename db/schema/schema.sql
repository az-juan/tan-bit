CREATE TYPE condicion AS ENUM ('nuevo', 'muy bueno', 'usado', 'reacondicionado', 'a reparar');
CREATE TYPE categoria AS ENUM ('notebook','pantalla','electronica','periferico','accesorio','audio','video','gaming','comunicacion','analogico','decoracion');
CREATE TYPE entrega   AS ENUM ('envio a domicilio','retiro en persona','ambos');

CREATE TABLE usuario (
    id_usuario  SERIAL      CONSTRAINT pk_usuario PRIMARY KEY,
    apellido    VARCHAR(50) NOT NULL,
    nombre      VARCHAR(50) NOT NULL,
    email       VARCHAR(50) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL
);

CREATE TABLE vendedor (
    id_vendedor     INTEGER  CONSTRAINT pk_vendedor PRIMARY KEY CONSTRAINT fk_vendedor_usuario REFERENCES usuario(id_usuario),
    telefono        BIGINT   UNIQUE NOT NULL,
    direccion_xmr   TEXT     UNIQUE     NULL,
    cbu             TEXT     UNIQUE NOT NULL,
    alias           TEXT     UNIQUE NOT NULL
);

CREATE TABLE articulo (
    id_articulo     SERIAL       CONSTRAINT pk_producto PRIMARY KEY,
    titulo          VARCHAR(50)  NOT NULL,
    descripcion     VARCHAR(250) NOT NULL,
    condicion       condicion    NOT NULL,
    categoria       categoria    NOT NULL,
    ruta_imagen     text         NULL,
    precio          DOUBLE PRECISION NOT NULL CHECK (precio >= 0),
    precio_xmr      DOUBLE PRECISION NOT NULL CHECK (precio_xmr >= 0),
    stock           INTEGER        NOT NULL CHECK (stock >= 0),
    fecha_publicacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ubicacion       text         NOT NULL,
    tipo_entrega    entrega      NOT NULL,
    costo_envio     DOUBLE PRECISION    NULL CHECK (costo_envio >= 0),
    zona_entrega    text             NULL,
    id_vendedor     INTEGER      NOT NULL CONSTRAINT fk_producto_vendedor REFERENCES vendedor(id_vendedor)
);

CREATE TABLE compra (
    id_compra       SERIAL          CONSTRAINT pk_compra PRIMARY KEY,
    id_comprador    INTEGER         NOT NULL CONSTRAINT fk_compra_usuario REFERENCES usuario(id_usuario),
    fecha           TIMESTAMP       NOT NULL,
    total           DOUBLE PRECISION  NOT NULL CHECK (total >= 1),
    total_xmr       DOUBLE PRECISION  NOT NULL CHECK (total >= 0)
);
CREATE TABLE detalle_compra (
    id_compra       INTEGER     NOT NULL CONSTRAINT fk_detalle_compra REFERENCES compra(id_compra),
    id_articulo     INTEGER     NOT NULL CONSTRAINT fk_detalle_producto REFERENCES articulo(id_articulo),
    cantidad        INTEGER     NOT NULL CHECK (cantidad >= 1),
    precio_historico_unitario DOUBLE PRECISION NOT NULL CHECK (precio_historico_unitario >= 1),
    precio_historico_unitario_xmr DOUBLE PRECISION NOT NULL CHECK (precio_historico_unitario_xmr >= 0),
    CONSTRAINT pk_detalle_compra PRIMARY KEY (id_compra, id_articulo)
);

CREATE TABLE carrito (
  id_usuario   int NOT NULL CONSTRAINT fk_usuario REFERENCES usuario(id_usuario),
  id_articulo  int NOT NULL CONSTRAINT fk_articulo REFERENCES articulo(id_articulo),
  cantidad     int NOT NULL CHECK (cantidad >= 1),
  fecha_agregado TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id_usuario, id_articulo)
);

CREATE TABLE favorito (
  id_usuario  int NOT NULL CONSTRAINT fk_usuario REFERENCES usuario(id_usuario),
  id_articulo int NOT NULL CONSTRAINT fk_articulo REFERENCES articulo(id_articulo),
  fecha       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id_usuario, id_articulo)
);
