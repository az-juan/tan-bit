CREATE TYPE condicion_enum AS ENUM ('nuevo', 'usado', 'reacondicionado', 'roto');

CREATE TABLE usuario (
    id_usuario SERIAL CONSTRAINT pk_usuario PRIMARY KEY,
    apellido VARCHAR(50),
    nombre VARCHAR(50),
    email VARCHAR(50) NOT NULL
);

CREATE TABLE vendedor (
    id_vendedor INTEGER CONSTRAINT pk_vendedor PRIMARY KEY CONSTRAINT fk_vendedor_usuario REFERENCES usuario(id_usuario),
    telefono VARCHAR(20) NOT NULL
);

CREATE TABLE producto (
    id_producto SERIAL CONSTRAINT pk_producto PRIMARY KEY,
    titulo VARCHAR(50) NOT NULL,
    descripcion VARCHAR(250) NOT NULL,
    condicion condicion_enum NOT NULL,
    precio NUMERIC(10, 2) NOT NULL,
    stock INTEGER NOT NULL,
    id_vendedor INTEGER NOT NULL CONSTRAINT fk_producto_vendedor REFERENCES vendedor(id_vendedor)
);

CREATE TABLE compra (
    id_compra SERIAL CONSTRAINT pk_compra PRIMARY KEY,
    id_comprador INTEGER NOT NULL CONSTRAINT fk_compra_usuario REFERENCES usuario(id_usuario),
    fecha TIMESTAMP NOT NULL,
    total NUMERIC(10, 2) NOT NULL
);

CREATE TABLE detalle_compra (
    id_compra INTEGER NOT NULL CONSTRAINT fk_detalle_compra REFERENCES compra(id_compra),
    id_producto INTEGER NOT NULL CONSTRAINT fk_detalle_producto REFERENCES producto(id_producto),
    cantidad INTEGER NOT NULL,
    precio_historico_unitario NUMERIC(10, 2) NOT NULL,
    CONSTRAINT pk_detalle_compra PRIMARY KEY (id_compra, id_producto)
);
