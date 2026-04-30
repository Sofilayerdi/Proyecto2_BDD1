package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "proy2-bck/db"
)

type Producto struct {
	IDProducto  int     `json:"id_producto"`
	Nombre      string  `json:"nombre"`
	Categoria   string  `json:"categoria"`
	IDProveedor int     `json:"id_proveedor"`
	Proveedor   string  `json:"proveedor,omitempty"`
	Precio      float64 `json:"precio"`
}

func ListarProductos(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 5
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	offset := limit * (page - 1)

	query := "SELECT p.id_producto, p.nombre, p.categoria, p.precio, pr.nombre AS proveedor 
				FROM producto p 
				LEFT JOIN proveedor pr ON p.id_proveedor = pr.id_proveedor
				LIMIT ? OFFSET ?;"

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		http.Error(w, "Error al obtener los productos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var productos []Productos
	for rows.Next() {
		var s Producto
		err := rows.Scan(&s.IDProducto, &s.Nombre, &s.Categoria, &s.Precio, &s.Proveedor)
		if err != nil {
			http.Error(w, "Error leyendo datos", http.StatusInternalServerError)
			return
		}
		productos = append(productos, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)

}

func ListarProductosCategoria(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 5
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	offset := limit * (page - 1)

	query := "SELECT p.id_producto, p.nombre, p.categoria, p.precio, pr.nombre AS proveedor 
				FROM producto p 
				LEFT JOIN proveedor pr ON p.id_proveedor = pr.id_proveedor
				LIMIT ? OFFSET ?;"

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		http.Error(w, "Error al obtener los productos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var productos []Productos
	for rows.Next() {
		var s Producto
		err := rows.Scan(&s.IDProducto, &s.Nombre, &s.Categoria, &s.Precio, &s.Proveedor)
		if err != nil {
			http.Error(w, "Error leyendo datos", http.StatusInternalServerError)
			return
		}
		productos = append(productos, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)

}

