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
	rows, err := db.DB.Query(`SELECT * FROM sp_reporte_ventas_mensuales()`)
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
	rows, err := db.DB.Query(`SELECT * FROM sp_top_productos()`)
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
