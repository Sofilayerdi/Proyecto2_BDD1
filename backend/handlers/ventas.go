package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	db "proy2-bck/db"

	"github.com/go-chi/chi/v5"
)

type Venta struct {
	IDVenta     int     `json:"id_venta"`
	IDCliente   int     `json:"id_cliente"`
	Cliente     string  `json:"cliente,omitempty"`
	IDEmpleado  int     `json:"id_empleado"`
	Empleado    string  `json:"empleado,omitempty"`
	Fecha       string  `json:"fecha"`
	PrecioTotal float64 `json:"precio_total"`
	Ramos       []int   `json:"ramos,omitempty"`
}

func ListarVentas(w http.ResponseWriter, r *http.Request) {
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
		SELECT v.id_venta, v.id_cliente, cl.nombre AS cliente,
		       v.id_empleado, em.nombre AS empleado,
		       v.fecha, v.precio_total
		FROM venta v
		JOIN cliente  cl ON v.id_cliente  = cl.id_cliente
		JOIN empleado em ON v.id_empleado = em.id_empleado
		ORDER BY v.fecha DESC
		LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		http.Error(w, "Error al obtener las ventas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ventas []Venta
	for rows.Next() {
		var v Venta
		err := rows.Scan(&v.IDVenta, &v.IDCliente, &v.Cliente,
			&v.IDEmpleado, &v.Empleado, &v.Fecha, &v.PrecioTotal)
		if err != nil {
			http.Error(w, "Error leyendo datos", http.StatusInternalServerError)
			return
		}
		ventas = append(ventas, v)
	}
	if ventas == nil {
		ventas = []Venta{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ventas)
}

func VerVenta(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var v Venta
	err = db.DB.QueryRow(`
		SELECT v.id_venta, v.id_cliente, cl.nombre AS cliente,
		       v.id_empleado, em.nombre AS empleado,
		       v.fecha, v.precio_total
		FROM venta v
		JOIN cliente  cl ON v.id_cliente  = cl.id_cliente
		JOIN empleado em ON v.id_empleado = em.id_empleado
		WHERE v.id_venta = $1`, id).
		Scan(&v.IDVenta, &v.IDCliente, &v.Cliente,
			&v.IDEmpleado, &v.Empleado, &v.Fecha, &v.PrecioTotal)

	if err == sql.ErrNoRows {
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error al obtener la venta", http.StatusInternalServerError)
		return
	}

	ramoRows, err := db.DB.Query(`SELECT id_ramo FROM ramo_venta WHERE id_venta=$1`, id)
	if err == nil {
		defer ramoRows.Close()
		for ramoRows.Next() {
			var rid int
			ramoRows.Scan(&rid)
			v.Ramos = append(v.Ramos, rid)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func CrearVenta(w http.ResponseWriter, r *http.Request) {
	var input Venta
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if input.IDCliente == 0 || input.IDEmpleado == 0 || input.Fecha == "" {
		http.Error(w, "id_cliente, id_empleado y fecha son requeridos", http.StatusBadRequest)
		return
	}
	if len(input.Ramos) == 0 {
		http.Error(w, "la venta debe incluir al menos un ramo", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Error iniciando transacción", http.StatusInternalServerError)
		return
	}

	var existeCliente bool
	tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM cliente WHERE id_cliente=$1)`,
		input.IDCliente).Scan(&existeCliente)
	if !existeCliente {
		tx.Rollback()
		http.Error(w, "Cliente no encontrado", http.StatusBadRequest)
		return
	}

	var existeEmpleado bool
	tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM empleado WHERE id_empleado=$1)`,
		input.IDEmpleado).Scan(&existeEmpleado)
	if !existeEmpleado {
		tx.Rollback()
		http.Error(w, "Empleado no encontrado", http.StatusBadRequest)
		return
	}

	var precioTotal float64
	for _, idRamo := range input.Ramos {
		var totalRamo float64
		err := tx.QueryRow(`SELECT total FROM ramo WHERE id_ramo=$1`, idRamo).Scan(&totalRamo)
		if err == sql.ErrNoRows {
			tx.Rollback()
			http.Error(w, "Ramo no encontrado: "+strconv.Itoa(idRamo), http.StatusBadRequest)
			return
		} else if err != nil {
			tx.Rollback()
			http.Error(w, "Error verificando ramo", http.StatusInternalServerError)
			return
		}
		precioTotal += totalRamo
	}

	var idVenta int
	err = tx.QueryRow(`
		INSERT INTO venta (id_cliente, id_empleado, fecha, precio_total)
		VALUES ($1, $2, $3, $4)
		RETURNING id_venta`,
		input.IDCliente, input.IDEmpleado, input.Fecha, precioTotal).
		Scan(&idVenta)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error creando venta", http.StatusInternalServerError)
		return
	}

	for _, idRamo := range input.Ramos {
		_, err = tx.Exec(`INSERT INTO ramo_venta (id_venta, id_ramo) VALUES ($1, $2)`,
			idVenta, idRamo)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error asociando ramo a la venta", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error confirmando transacción", http.StatusInternalServerError)
		return
	}

	input.IDVenta = idVenta
	input.PrecioTotal = precioTotal
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func EliminarVenta(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	res, err := db.DB.Exec(`DELETE FROM venta WHERE id_venta=$1`, id)
	if err != nil {
		http.Error(w, "Error al eliminar la venta", http.StatusInternalServerError)
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
