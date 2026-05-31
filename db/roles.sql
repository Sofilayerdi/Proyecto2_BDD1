CREATE ROLE superadmin;
CREATE ROLE gerente;
CREATE ROLE vendedor;
CREATE ROLE comprador;
CREATE ROLE auditor;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO superadmin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO superadmin;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO gerente;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO gerente;

GRANT SELECT, INSERT, UPDATE ON producto, ramo, ramo_producto, venta, ramo_venta TO vendedor;
GRANT SELECT ON proveedor, cliente TO vendedor;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO vendedor;

GRANT SELECT ON venta, ramo_venta, ramo_producto, producto, ramo, cliente, empleado, proveedor TO auditor;

GRANT SELECT ON producto, cliente, empleado TO comprador;
GRANT SELECT, INSERT ON ramo, ramo_producto, venta, ramo_venta TO comprador;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO comprador;

CREATE USER superadmin1 WITH PASSWORD 'secret' IN ROLE superadmin;
CREATE USER gerente1 WITH PASSWORD 'secret' IN ROLE gerente;
CREATE USER vendedor1 WITH PASSWORD 'secret' IN ROLE vendedor;
CREATE USER auditor1 WITH PASSWORD 'secret' IN ROLE auditor;
CREATE USER comprador1 WITH PASSWORD 'secret' IN ROLE comprador;