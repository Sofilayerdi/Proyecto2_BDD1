-- SL 1
CREATE OR REPLACE PROCEDURE sp_crear_venta(
    IN  p_id_cliente   INT,
    IN  p_id_empleado  INT,
    IN  p_fecha        DATE,
    IN  p_ramos        INT[],
    OUT p_id_venta     INT,
    OUT p_precio_total DECIMAL(10,2),
    OUT p_mensaje      VARCHAR(200)
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_ramo_id    INT;
    v_total_ramo DECIMAL(10,2);
BEGIN
    IF NOT EXISTS (SELECT 1 FROM cliente WHERE id_cliente = p_id_cliente) THEN
        p_mensaje := 'Cliente no encontrado';
        ROLLBACK;
        RETURN;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM empleado WHERE id_empleado = p_id_empleado) THEN
        p_mensaje := 'Empleado no encontrado';
        ROLLBACK;
        RETURN;
    END IF;

    p_precio_total := 0;
    FOREACH v_ramo_id IN ARRAY p_ramos LOOP
        SELECT total INTO v_total_ramo FROM ramo WHERE id_ramo = v_ramo_id;
        IF NOT FOUND THEN
            p_mensaje := 'Ramo no encontrado: ' || v_ramo_id;
            ROLLBACK;
            RETURN;
        END IF;
        p_precio_total := p_precio_total + v_total_ramo;
    END LOOP;

    INSERT INTO venta (id_cliente, id_empleado, fecha, precio_total)
    VALUES (p_id_cliente, p_id_empleado, p_fecha, p_precio_total)
    RETURNING id_venta INTO p_id_venta;

    FOREACH v_ramo_id IN ARRAY p_ramos LOOP
        INSERT INTO ramo_venta (id_venta, id_ramo)
        VALUES (p_id_venta, v_ramo_id);
    END LOOP;

    p_mensaje := 'Venta creada exitosamente';

EXCEPTION
    WHEN OTHERS THEN
        p_mensaje := 'Error inesperado: ' || SQLERRM;
        ROLLBACK;
END;
$$;

-- SP 2
CREATE OR REPLACE PROCEDURE sp_crear_ramo(
    IN  p_productos    INT[],
    IN  p_cantidades   INT[],
    OUT p_id_ramo      INT,
    OUT p_total        DECIMAL(10,2),
    OUT p_mensaje      VARCHAR(200)
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_i          INT;
    v_precio     DECIMAL(10,2);
    v_stock      INT;
    v_id_prod    INT;
    v_cantidad   INT;
BEGIN
    p_total := 0;

    FOR v_i IN 1..array_length(p_productos, 1) LOOP
        v_id_prod  := p_productos[v_i];
        v_cantidad := p_cantidades[v_i];

        SELECT precio, cantidad INTO v_precio, v_stock
        FROM producto WHERE id_producto = v_id_prod;

        IF NOT FOUND THEN
            p_mensaje := 'Producto no encontrado: ' || v_id_prod;
            ROLLBACK;
            RETURN;
        END IF;

        IF v_stock < v_cantidad THEN
            p_mensaje := 'Stock insuficiente para producto: ' || v_id_prod;
            ROLLBACK;
            RETURN;
        END IF;

        p_total := p_total + (v_precio * v_cantidad);
    END LOOP;

    INSERT INTO ramo (total) VALUES (p_total)
    RETURNING id_ramo INTO p_id_ramo;

    FOR v_i IN 1..array_length(p_productos, 1) LOOP
        v_id_prod  := p_productos[v_i];
        v_cantidad := p_cantidades[v_i];

        INSERT INTO ramo_producto (id_ramo, id_producto, cantidad)
        VALUES (p_id_ramo, v_id_prod, v_cantidad);

        UPDATE producto
        SET cantidad = cantidad - v_cantidad
        WHERE id_producto = v_id_prod;
    END LOOP;

    p_mensaje := 'Ramo creado exitosamente';

EXCEPTION
    WHEN OTHERS THEN
        p_mensaje := 'Error inesperado: ' || SQLERRM;
        ROLLBACK;
END;
$$;

-- SP 3
CREATE OR REPLACE PROCEDURE sp_actualizar_stock(
    IN  p_id_producto INT,
    IN  p_cantidad    INT,
    OUT p_stock_nuevo INT,
    OUT p_mensaje     VARCHAR(200)
)
LANGUAGE plpgsql
AS $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM producto WHERE id_producto = p_id_producto) THEN
        p_mensaje := 'Producto no encontrado';
        RETURN;
    END IF;

    IF p_cantidad < 0 THEN
        SELECT cantidad INTO p_stock_nuevo FROM producto WHERE id_producto = p_id_producto;
        IF p_stock_nuevo + p_cantidad < 0 THEN
            p_mensaje := 'Stock insuficiente';
            RETURN;
        END IF;
    END IF;

    UPDATE producto
    SET cantidad = cantidad + p_cantidad
    WHERE id_producto = p_id_producto
    RETURNING cantidad INTO p_stock_nuevo;

    p_mensaje := 'Stock actualizado correctamente';

EXCEPTION
    WHEN OTHERS THEN
        p_mensaje := 'Error inesperado: ' || SQLERRM;
END;
$$;

-- SP 4
CREATE OR REPLACE FUNCTION sp_reporte_ventas_mensuales()
RETURNS TABLE (
    mes          TEXT,
    total_ventas BIGINT,
    ingresos     DECIMAL(10,2)
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        TO_CHAR(fecha, 'YYYY-MM') AS mes,
        COUNT(*)::BIGINT           AS total_ventas,
        SUM(precio_total)          AS ingresos
    FROM venta
    GROUP BY TO_CHAR(fecha, 'YYYY-MM')
    ORDER BY mes ASC;

EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error en reporte: %', SQLERRM;
END;
$$;

-- SP 5
CREATE OR REPLACE FUNCTION sp_top_productos()
RETURNS TABLE (
    producto      TEXT,
    categoria     TEXT,
    proveedor     TEXT,
    total_vendido BIGINT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        p.nombre::TEXT       AS producto,
        p.categoria::TEXT    AS categoria,
        pr.nombre::TEXT      AS proveedor,
        SUM(rp.cantidad)::BIGINT AS total_vendido
    FROM ramo_producto rp
    JOIN producto  p  ON rp.id_producto = p.id_producto
    JOIN proveedor pr ON p.id_proveedor = pr.id_proveedor
    JOIN ramo_venta rv ON rp.id_ramo    = rv.id_ramo
    GROUP BY p.id_producto, p.nombre, p.categoria, pr.nombre
    ORDER BY total_vendido DESC
    LIMIT 20;

EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error en reporte: %', SQLERRM;
END;
$$;