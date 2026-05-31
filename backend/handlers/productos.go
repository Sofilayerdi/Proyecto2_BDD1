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

// crear producto con ORM
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

	modelo := db.Producto{
		Nombre:      p.Nombre,
		Categoria:   p.Categoria,
		IDProveedor: p.IDProveedor,
		Cantidad:    p.Cantidad,
		Precio:      p.Precio,
	}

	if result := db.GORM.Create(&modelo); result.Error != nil {
		http.Error(w, "Error al crear producto: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	p.IDProducto = modelo.IDProducto
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)

}

// editar con ORM
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

	var modelo db.Producto
	if result := db.GORM.First(&modelo, id); result.Error != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	if result := db.GORM.Model(&modelo).Updates(db.Producto{
		Nombre:      p.Nombre,
		Categoria:   p.Categoria,
		IDProveedor: p.IDProveedor,
		Cantidad:    p.Cantidad,
		Precio:      p.Precio,
	}); result.Error != nil {
		http.Error(w, "Error al editar producto", http.StatusInternalServerError)
		return
	}

	p.IDProducto = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// eliminar con ORM
func EliminarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var modelo db.Producto
	if result := db.GORM.First(&modelo, id); result.Error != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	if result := db.GORM.Delete(&modelo); result.Error != nil {
		http.Error(w, "Error al eliminar producto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
