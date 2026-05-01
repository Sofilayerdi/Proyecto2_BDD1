package handlers

import (
	"encoding/json"
	"net/http"

	db "proy2-bck/db"
)

type VentaPorEmpleado struct {
	Empleado    string  `json:"empleado"`
	Rol         string  `json:"rol"`
	TotalVentas int     `json:"total_ventas"`
	Ingresos    float64 `json:"ingresos"`
}

type ProductoEnRamo struct {
	Producto     string `json:"producto"`
	Categoria    string `json:"categoria"`
	Proveedor    string `json:"proveedor"`
	TotalVendido int    `json:"total_vendido"`
}

type ItemInventario struct {
	IDProducto int     `json:"id_producto"`
	Nombre     string  `json:"nombre"`
	Categoria  string  `json:"categoria"`
	Proveedor  string  `json:"proveedor"`
	Cantidad   int     `json:"cantidad"`
	Precio     float64 `json:"precio"`
}

type VistaVenta struct {
	IDVenta     int     `json:"id_venta"`
	Fecha       string  `json:"fecha"`
	Cliente     string  `json:"cliente"`
	Empleado    string  `json:"empleado"`
	PrecioTotal float64 `json:"precio_total"`
}

func VentasPorEmpleado(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT empleado, rol, total_ventas, ingresos
		FROM vista_ventas_empleado
		ORDER BY ingresos DESC`)
	if err != nil {
		http.Error(w, "Error en reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []VentaPorEmpleado
	for rows.Next() {
		var v VentaPorEmpleado
		rows.Scan(&v.Empleado, &v.Rol, &v.TotalVentas, &v.Ingresos)
		lista = append(lista, v)
	}
	if lista == nil {
		lista = []VentaPorEmpleado{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func ProductosEnRamos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		WITH vendidos AS (
			SELECT rp.id_producto, SUM(rp.cantidad) AS total_vendido
			FROM ramo_producto rp
			JOIN ramo_venta rv ON rp.id_ramo = rv.id_ramo
			GROUP BY rp.id_producto
		)
		SELECT p.nombre    AS producto,
		       p.categoria AS categoria,
		       pr.nombre   AS proveedor,
		       v.total_vendido
		FROM vendidos v
		JOIN producto  p  ON v.id_producto = p.id_producto
		JOIN proveedor pr ON p.id_proveedor = pr.id_proveedor
		ORDER BY v.total_vendido DESC`)
	if err != nil {
		http.Error(w, "Error en reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []ProductoEnRamo
	for rows.Next() {
		var p ProductoEnRamo
		rows.Scan(&p.Producto, &p.Categoria, &p.Proveedor, &p.TotalVendido)
		lista = append(lista, p)
	}
	if lista == nil {
		lista = []ProductoEnRamo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func Inventario(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT p.id_producto,
		       p.nombre,
		       p.categoria,
		       pr.nombre AS proveedor,
		       p.cantidad,
		       p.precio
		FROM producto p
		JOIN proveedor pr ON p.id_proveedor = pr.id_proveedor
		WHERE EXISTS (
		    SELECT 1
		    FROM (SELECT AVG(cantidad) AS promedio FROM producto) avg_t
		    WHERE p.cantidad < avg_t.promedio
		)
		ORDER BY p.cantidad ASC`)
	if err != nil {
		http.Error(w, "Error en reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []ItemInventario
	for rows.Next() {
		var p ItemInventario
		rows.Scan(&p.IDProducto, &p.Nombre, &p.Categoria, &p.Proveedor, &p.Cantidad, &p.Precio)
		lista = append(lista, p)
	}
	if lista == nil {
		lista = []ItemInventario{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func VistaVentas(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT id_venta, fecha, cliente, empleado, precio_total
		FROM vista_detalle_ventas
		ORDER BY fecha DESC`)
	if err != nil {
		http.Error(w, "Error en reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []VistaVenta
	for rows.Next() {
		var v VistaVenta
		rows.Scan(&v.IDVenta, &v.Fecha, &v.Cliente, &v.Empleado, &v.PrecioTotal)
		lista = append(lista, v)
	}
	if lista == nil {
		lista = []VistaVenta{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}
