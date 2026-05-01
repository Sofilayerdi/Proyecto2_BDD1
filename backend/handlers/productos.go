package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	db "proy2-bck/db"

	"github.com/go-chi/chi/v5"
)

type Producto struct {
	IDProducto  int     `json:"id_producto"`
	Nombre      string  `json:"nombre"`
	Categoria   string  `json:"categoria"`
	IDProveedor int     `json:"id_proveedor"`
	Proveedor   string  `json:"proveedor,omitempty"`
	Cantidad    int     `json:"cantidad"`
	Precio      float64 `json:"precio"`
}

func ListarProductos(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	offset := limit * (page - 1)
	cat := r.URL.Query().Get("categoria")

	var rows *sql.Rows
	var err error

	if cat != "" {
		rows, err = db.DB.Query(`
			SELECT p.id_producto, p.nombre, p.categoria,
			       p.id_proveedor, pr.nombre AS proveedor,
			       p.cantidad, p.precio
			FROM producto p
			JOIN proveedor pr ON p.id_proveedor = pr.id_proveedor
			WHERE p.categoria = $1
			ORDER BY p.nombre
			LIMIT $2 OFFSET $3`,
			cat, limit, offset)
	} else {
		rows, err = db.DB.Query(`
			SELECT p.id_producto, p.nombre, p.categoria,
			       p.id_proveedor, pr.nombre AS proveedor,
			       p.cantidad, p.precio
			FROM producto p
			JOIN proveedor pr ON p.id_proveedor = pr.id_proveedor
			ORDER BY p.categoria, p.nombre
			LIMIT $1 OFFSET $2`,
			limit, offset)
	}
	if err != nil {
		http.Error(w, "Error al obtener los productos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var productos []Producto
	for rows.Next() {
		var p Producto
		err := rows.Scan(&p.IDProducto, &p.Nombre, &p.Categoria,
			&p.IDProveedor, &p.Proveedor, &p.Cantidad, &p.Precio)
		if err != nil {
			http.Error(w, "Error leyendo datos", http.StatusInternalServerError)
			return
		}
		productos = append(productos, p)
	}
	if productos == nil {
		productos = []Producto{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

func VerProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p Producto
	err = db.DB.QueryRow(`
		SELECT p.id_producto, p.nombre, p.categoria,
		       p.id_proveedor, pr.nombre AS proveedor,
		       p.cantidad, p.precio
		FROM producto p
		JOIN proveedor pr ON p.id_proveedor = pr.id_proveedor
		WHERE p.id_producto = $1`, id).
		Scan(&p.IDProducto, &p.Nombre, &p.Categoria,
			&p.IDProveedor, &p.Proveedor, &p.Cantidad, &p.Precio)

	if err == sql.ErrNoRows {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error al obtener el producto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func CrearProducto(w http.ResponseWriter, r *http.Request) {
	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if p.Nombre == "" || p.Categoria == "" || p.IDProveedor == 0 {
		http.Error(w, "nombre, categoria e id_proveedor son requeridos", http.StatusBadRequest)
		return
	}
	categorias := map[string]bool{"flor": true, "follaje": true, "liston": true, "papel": true}
	if !categorias[p.Categoria] {
		http.Error(w, "categoria debe ser: flor, follaje, liston o papel", http.StatusBadRequest)
		return
	}
	if p.Precio < 0 || p.Cantidad < 0 {
		http.Error(w, "precio y cantidad no pueden ser negativos", http.StatusBadRequest)
		return
	}

	err := db.DB.QueryRow(`
		INSERT INTO producto (nombre, categoria, id_proveedor, cantidad, precio)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_producto`,
		p.Nombre, p.Categoria, p.IDProveedor, p.Cantidad, p.Precio).
		Scan(&p.IDProducto)
	if err != nil {
		http.Error(w, "Error al crear producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func EditarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if p.Nombre == "" || p.Categoria == "" || p.IDProveedor == 0 {
		http.Error(w, "nombre, categoria e id_proveedor son requeridos", http.StatusBadRequest)
		return
	}
	categorias := map[string]bool{"flor": true, "follaje": true, "liston": true, "papel": true}
	if !categorias[p.Categoria] {
		http.Error(w, "categoria debe ser: flor, follaje, liston o papel", http.StatusBadRequest)
		return
	}
	if p.Precio < 0 || p.Cantidad < 0 {
		http.Error(w, "precio y cantidad no pueden ser negativos", http.StatusBadRequest)
		return
	}

	res, err := db.DB.Exec(`
		UPDATE producto
		SET nombre=$1, categoria=$2, id_proveedor=$3, cantidad=$4, precio=$5
		WHERE id_producto=$6`,
		p.Nombre, p.Categoria, p.IDProveedor, p.Cantidad, p.Precio, id)
	if err != nil {
		http.Error(w, "Error al editar producto", http.StatusInternalServerError)
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	p.IDProducto = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func EliminarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	res, err := db.DB.Exec(`DELETE FROM producto WHERE id_producto=$1`, id)
	if err != nil {
		http.Error(w, "No se puede eliminar: el producto está en uso", http.StatusConflict)
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
