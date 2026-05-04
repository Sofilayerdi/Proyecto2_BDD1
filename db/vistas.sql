CREATE VIEW vista_ventas_empleado AS
SELECT
    e.nombre        AS empleado,
    e.rol           AS rol,
    COUNT(v.id_venta)       AS total_ventas,
    SUM(v.precio_total)     AS ingresos
FROM empleado e
LEFT JOIN venta v ON e.id_empleado = v.id_empleado
GROUP BY e.id_empleado, e.nombre, e.rol
HAVING COUNT(v.id_venta) > 0
ORDER BY ingresos DESC;

CREATE VIEW vista_detalle_ventas AS
SELECT
    v.id_venta,
    v.fecha::TEXT           AS fecha,
    c.nombre                AS cliente,
    e.nombre                AS empleado,
    v.precio_total
FROM venta v
JOIN cliente  c ON v.id_cliente  = c.id_cliente
JOIN empleado e ON v.id_empleado = e.id_empleado;