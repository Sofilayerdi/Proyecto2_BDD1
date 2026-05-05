package handlers

import (
	"encoding/json"
	"net/http"

	db "proy2-bck/db"
)

type VentaMensual struct {
	Mes         string  `json:"mes"`
	TotalVentas int     `json:"total_ventas"`
	Ingresos    float64 `json:"ingresos"`
}

type ProductoVendido struct {
	Producto     string `json:"producto"`
	Categoria    string `json:"categoria"`
	Proveedor    string `json:"proveedor"`
	TotalVendido int    `json:"total_vendido"`
}

func VentasMensuales(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT
			TO_CHAR(fecha, 'YYYY-MM')   AS mes,
			COUNT(*)                     AS total_ventas,
			SUM(precio_total)            AS ingresos
		FROM venta
		GROUP BY TO_CHAR(fecha, 'YYYY-MM')
		ORDER BY mes ASC`)
	if err != nil {
		http.Error(w, "Error en reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []VentaMensual
	for rows.Next() {
		var v VentaMensual
		rows.Scan(&v.Mes, &v.TotalVentas, &v.Ingresos)
		lista = append(lista, v)
	}
	if lista == nil {
		lista = []VentaMensual{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func TopProductosVendidos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT
			p.nombre       AS producto,
			p.categoria    AS categoria,
			pr.nombre      AS proveedor,
			SUM(rp.cantidad) AS total_vendido
		FROM ramo_producto rp
		JOIN producto  p  ON rp.id_producto  = p.id_producto
		JOIN proveedor pr ON p.id_proveedor   = pr.id_proveedor
		JOIN ramo_venta rv ON rp.id_ramo      = rv.id_ramo
		GROUP BY p.id_producto, p.nombre, p.categoria, pr.nombre
		ORDER BY total_vendido DESC
		LIMIT 20`)
	if err != nil {
		http.Error(w, "Error en reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []ProductoVendido
	for rows.Next() {
		var p ProductoVendido
		rows.Scan(&p.Producto, &p.Categoria, &p.Proveedor, &p.TotalVendido)
		lista = append(lista, p)
	}
	if lista == nil {
		lista = []ProductoVendido{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}
