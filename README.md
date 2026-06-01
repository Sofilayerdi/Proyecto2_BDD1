# Proyecto 3 BDD1 

---

## Tecnologías utilizadas
- **Frontend:** React
- **Backend:** Go
- **Base de datos:** PostgreSQL
- **ORM:** GORM
- **Contenedores:** Docker + Docker Compose

---

## Requisitos previos

Antes de correr el proyecto, asegurate de tener instalado:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

## Configuración inicial

1. **Clonar el repositorio y cambiar a la rama correcta:**

```bash
git clone https://github.com/Sofilayerdi/Proyecto2_BDD1.git
cd Proyecto2_BDD1
git checkout proyecto-3
```

2. **Crear el archivo de variables de entorno:**

Copiar el archivo de ejemplo y ajustar si es necesario:

```bash
cp .env.example .env
```

El archivo `.env` debe contener:

```env
DB_USER=proy3
DB_PASSWORD=secret
DB_NAME=proyecto2
DB_HOST=db
DB_PORT=5432
---

## Correr el proyecto

Desde la raíz del proyecto, ejecutar:

```bash
docker compose up --build
```

Esto levantará tres servicios:

| Servicio   | Puerto local |
|------------|-------------|
| `db`       | `5435`      |
| `backend`  | `8000`      |
| `frontend` | `5173`      |

Una vez que todos los contenedores estén corriendo, abrir el navegador en:

```
http://localhost:5173
```

---

## Usuarios de prueba

Todos los usuarios tienen la contraseña: `secret`

| Usuario      | Rol        | 
|--------------|------------|
| superadmin1  | superadmin |
| gerente1     | gerente    | 
| vendedor1    | vendedor   |
| auditor1     | auditor    |
| comprador1   | comprador  |

---

## Esquema de roles

### 1. superadmin
Acceso total al sistema sin restricciones. Todas las operaciones (SELECT, INSERT, UPDATE, DELETE) en todas las tablas.

**UI:** Inventario completo · Compras · Reportes

---

### 2. gerente
Acceso operativo completo incluyendo reportes.

**Tablas:** producto, ramo, ramo_producto, venta, ramo_venta, cliente, empleado, proveedor, usuario — SELECT, INSERT, UPDATE, DELETE

**UI:** Inventario completo · Compras · Reportes

---

### 3. vendedor
Gestión de inventario y ventas. Sin acceso a eliminar registros, tabla de usuarios ni reportes.

**Tablas:** producto, ramo, ramo_producto, venta, ramo_venta — SELECT, INSERT, UPDATE · cliente, empleado, proveedor — SELECT · usuario — sin acceso

**UI:** Inventario (ver, crear, editar) · Compras · Sin reportes

---

### 4. auditor
Solo lectura. Acceso exclusivo para generar reportes.

**Tablas:** producto, ramo, ramo_producto, venta, ramo_venta, cliente, empleado, proveedor — SELECT · usuario — sin acceso

**UI:** Solo reportes

---

### 5. comprador
Puede ver productos y realizar compras. Sin acceso a gestión de inventario ni reportes.

**Tablas:** producto, ramo, ramo_producto, venta, ramo_venta, cliente, empleado — SELECT/INSERT según aplique · proveedor, usuario — sin acceso

**UI:** Inventario (solo ver) · Compras · Sin reportes

---

## Stored Procedures

El backend invoca los siguientes stored procedures desde Go:

1. **`sp_crear_venta`** *(PROCEDURE)* — Registra una venta completa. Recibe id de cliente, empleado, fecha y un arreglo de ramos. Valida que el cliente y empleado existan, calcula el precio total sumando los ramos y los vincula a la venta. Incluye parámetros de salida (`p_id_venta`, `p_precio_total`, `p_mensaje`), manejo de excepciones y `ROLLBACK` explícito ante cualquier error.

2. **`sp_crear_ramo`** *(PROCEDURE)* — Crea un ramo a partir de un arreglo de productos y cantidades. Valida stock disponible por producto, calcula el total, inserta el ramo, registra el detalle en `ramo_producto` y descuenta el stock. Incluye parámetros de salida (`p_id_ramo`, `p_total`, `p_mensaje`) y `ROLLBACK` ante stock insuficiente o producto no encontrado.

3. **`sp_actualizar_stock`** *(PROCEDURE)* — Actualiza el stock de un producto sumando o restando una cantidad. Valida que el producto exista y que el stock no quede negativo. Retorna el stock nuevo y un mensaje de resultado mediante parámetros de salida.

4. **`sp_reporte_ventas_mensuales`** *(FUNCTION)* — Devuelve un reporte agrupado por mes con el total de ventas e ingresos. Incluye manejo de excepciones con `RAISE EXCEPTION`.

5. **`sp_top_productos`** *(FUNCTION)* — Devuelve los 20 productos más vendidos con nombre, categoría, proveedor y unidades totales vendidas, cruzando `ramo_producto`, `producto`, `proveedor` y `ramo_venta`. Incluye manejo de excepciones.

> `sp_crear_venta` y `sp_crear_ramo` implementan transacciones explícitas con `ROLLBACK`, parámetros de entrada/salida y manejo de excepciones con `WHEN OTHERS`.

---

## ORM

Se utiliza GORM para las siguientes operaciones CRUD:

- Consulta de productos (`SELECT`)
- Creación y edición de productos (`INSERT` / `UPDATE`)
- Consulta de clientes y empleados (`SELECT`)

---

## Detener el proyecto

```bash
docker compose down
```

Para detener **y eliminar los volúmenes** (borra los datos de la base de datos):

```bash
docker compose down -v
```