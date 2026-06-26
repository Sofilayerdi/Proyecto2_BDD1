CREATE TABLE proveedor (
    id_proveedor SERIAL PRIMARY KEY,
    nombre       VARCHAR(100) NOT NULL,
    correo       VARCHAR(100) NOT NULL,
    telefono     VARCHAR(20)  NOT NULL
);

CREATE TABLE producto (
    id_producto  SERIAL PRIMARY KEY,
    nombre       VARCHAR(100)  NOT NULL,
    categoria    VARCHAR(20)   NOT NULL,
    id_proveedor INT           NOT NULL,
    cantidad     INT           NOT NULL,
    precio       DECIMAL(10,2) NOT NULL,

    CONSTRAINT fk_producto_proveedor
        FOREIGN KEY (id_proveedor)
        REFERENCES proveedor(id_proveedor)
        ON DELETE RESTRICT ON UPDATE CASCADE,

    CONSTRAINT chk_categoria_producto
        CHECK (categoria IN ('flor', 'follaje', 'liston', 'papel')),

    CONSTRAINT chk_cantidad_producto
        CHECK (cantidad >= 0),

    CONSTRAINT chk_precio_producto
        CHECK (precio >= 0)
);

CREATE TABLE ramo (
    id_ramo SERIAL PRIMARY KEY,
    total   DECIMAL(10,2) NOT NULL,

    CONSTRAINT chk_total_ramo
        CHECK (total >= 0)
);

CREATE TABLE ramo_producto (
    id_ramo     INT NOT NULL,
    id_producto INT NOT NULL,
    cantidad    INT NOT NULL,

    PRIMARY KEY (id_ramo, id_producto),

    CONSTRAINT fk_rp_ramo
        FOREIGN KEY (id_ramo)
        REFERENCES ramo(id_ramo)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_rp_producto
        FOREIGN KEY (id_producto)
        REFERENCES producto(id_producto)
        ON DELETE RESTRICT ON UPDATE CASCADE,

    CONSTRAINT chk_cantidad_rp
        CHECK (cantidad > 0)
);

CREATE TABLE cliente (
    id_cliente SERIAL PRIMARY KEY,
    nombre     VARCHAR(100) NOT NULL,
    telefono   VARCHAR(20)  NOT NULL,
    correo     VARCHAR(100) NOT NULL,
    direccion  VARCHAR(200) NOT NULL
);

CREATE TABLE empleado (
    id_empleado SERIAL PRIMARY KEY,
    nombre      VARCHAR(100) NOT NULL,
    rol         VARCHAR(50)  NOT NULL,
    telefono    VARCHAR(20)  NOT NULL,
    correo      VARCHAR(100) NOT NULL
);

CREATE TABLE usuario (
    id_usuario  SERIAL PRIMARY KEY,
    username    VARCHAR(50) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    rol         VARCHAR(20) NOT NULL CHECK (rol IN ('superadmin','gerente','vendedor','auditor','comprador')),
    id_empleado INT REFERENCES empleado(id_empleado) ON DELETE SET NULL,
    id_cliente  INT REFERENCES cliente(id_cliente) ON DELETE SET NULL
);

CREATE TABLE venta (
    id_venta     SERIAL PRIMARY KEY,
    id_cliente   INT           NOT NULL,
    id_empleado  INT           NOT NULL,
    fecha        DATE          NOT NULL,
    precio_total DECIMAL(10,2) NOT NULL,

    CONSTRAINT fk_venta_cliente
        FOREIGN KEY (id_cliente)
        REFERENCES cliente(id_cliente)
        ON DELETE RESTRICT ON UPDATE CASCADE,

    CONSTRAINT fk_venta_empleado
        FOREIGN KEY (id_empleado)
        REFERENCES empleado(id_empleado)
        ON DELETE RESTRICT ON UPDATE CASCADE,

    CONSTRAINT chk_precio_venta
        CHECK (precio_total >= 0)
);

CREATE TABLE ramo_venta (
    id_venta INT NOT NULL,
    id_ramo  INT NOT NULL,

    PRIMARY KEY (id_venta, id_ramo),

    CONSTRAINT fk_rv_venta
        FOREIGN KEY (id_venta)
        REFERENCES venta(id_venta)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_rv_ramo
        FOREIGN KEY (id_ramo)
        REFERENCES ramo(id_ramo)
        ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_venta_fecha      ON venta(fecha);
CREATE INDEX idx_ramoprod_prod    ON ramo_producto(id_producto);
