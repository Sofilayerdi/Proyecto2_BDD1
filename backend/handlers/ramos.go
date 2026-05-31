package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	db "proy2-bck/db"

	"github.com/go-chi/chi/v5"
)

type RamoProductoItem struct {
	IDProducto int     `json:"id_producto"`
	Nombre     string  `json:"nombre,omitempty"`
	Cantidad   int     `json:"cantidad"`
	Precio     float64 `json:"precio,omitempty"`
}

type Ramo struct {
	IDRamo    int                `json:"id_ramo"`
	Total     float64            `json:"total"`
	Productos []RamoProductoItem `json:"productos,omitempty"`
}

func ListarRamos(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	offset := limit * (page - 1)

	rows, err := db.DB.Query(`
		SELECT id_ramo, total
		FROM ramo
		ORDER BY id_ramo DESC
		LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		http.Error(w, "Error al obtener ramos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ramos []Ramo
	for rows.Next() {
		var rm Ramo
		if err := rows.Scan(&rm.IDRamo, &rm.Total); err != nil {
			http.Error(w, "Error leyendo datos", http.StatusInternalServerError)
			return
		}
		ramos = append(ramos, rm)
	}
	if ramos == nil {
		ramos = []Ramo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ramos)
}

func VerRamo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var rm Ramo
	err = db.DB.QueryRow(`SELECT id_ramo, total FROM ramo WHERE id_ramo=$1`, id).
		Scan(&rm.IDRamo, &rm.Total)
	if err == sql.ErrNoRows {
		http.Error(w, "Ramo no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error al obtener el ramo", http.StatusInternalServerError)
		return
	}

	prodRows, err := db.DB.Query(`
		SELECT rp.id_producto, p.nombre, rp.cantidad, p.precio
		FROM ramo_producto rp
		JOIN producto p ON rp.id_producto = p.id_producto
		WHERE rp.id_ramo = $1`, id)
	if err != nil {
		http.Error(w, "Error al obtener productos del ramo", http.StatusInternalServerError)
		return
	}
	defer prodRows.Close()

	for prodRows.Next() {
		var item RamoProductoItem
		prodRows.Scan(&item.IDProducto, &item.Nombre, &item.Cantidad, &item.Precio)
		rm.Productos = append(rm.Productos, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rm)
}

// Crear ramo con GORM
func CrearRamo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Productos []RamoProductoItem `json:"productos"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if len(input.Productos) == 0 {
		http.Error(w, "el ramo debe tener al menos un producto", http.StatusBadRequest)
		return
	}

	tx := db.GORM.Begin()
	if tx.Error != nil {
		http.Error(w, "Error iniciando transacción", http.StatusInternalServerError)
		return
	}

	var total float64
	for _, item := range input.Productos {
		if item.Cantidad <= 0 {
			tx.Rollback()
			http.Error(w, "La cantidad debe ser mayor a 0", http.StatusBadRequest)
			return
		}

		var producto db.Producto
		if result := tx.First(&producto, item.IDProducto); result.Error != nil {
			tx.Rollback()
			http.Error(w, "Producto no encontrado: "+strconv.Itoa(item.IDProducto), http.StatusBadRequest)
			return
		}

		if producto.Cantidad < item.Cantidad {
			tx.Rollback()
			http.Error(w,
				"Stock insuficiente para producto "+strconv.Itoa(item.IDProducto)+
					" (disponible: "+strconv.Itoa(producto.Cantidad)+
					", solicitado: "+strconv.Itoa(item.Cantidad)+")",
				http.StatusConflict)
			return
		}

		total += producto.Precio * float64(item.Cantidad)
	}

	nuevoRamo := db.Ramo{Total: total}
	if result := tx.Create(&nuevoRamo); result.Error != nil {
		tx.Rollback()
		http.Error(w, "Error creando el ramo", http.StatusInternalServerError)
		return
	}

	for _, item := range input.Productos {
		ramoProducto := db.RamoProducto{
			IDRamo:     nuevoRamo.IDRamo,
			IDProducto: item.IDProducto,
			Cantidad:   item.Cantidad,
		}
		if result := tx.Create(&ramoProducto); result.Error != nil {
			tx.Rollback()
			http.Error(w, "Error insertando producto en ramo", http.StatusInternalServerError)
			return
		}

		if result := tx.Model(&db.Producto{}).
			Where("id_producto = ?", item.IDProducto).
			Update("cantidad", db.GORM.Raw("cantidad - ?", item.Cantidad)); result.Error != nil {
			tx.Rollback()
			http.Error(w, "Error actualizando stock", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error confirmando transacción", http.StatusInternalServerError)
		return
	}

	resp := Ramo{IDRamo: nuevoRamo.IDRamo, Total: total, Productos: input.Productos}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func EliminarRamo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var enVenta bool
	db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM ramo_venta WHERE id_ramo=$1)`, id).
		Scan(&enVenta)
	if enVenta {
		http.Error(w, "No se puede eliminar: el ramo ya está asociado a una venta", http.StatusConflict)
		return
	}

	res, err := db.DB.Exec(`DELETE FROM ramo WHERE id_ramo=$1`, id)
	if err != nil {
		http.Error(w, "Error al eliminar el ramo", http.StatusInternalServerError)
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "Ramo no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
