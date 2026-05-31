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

	var idVenta int
	var precioTotal float64
	var mensaje string

	err := db.DB.QueryRow(
		`CALL sp_crear_venta($1, $2, $3, $4, $5, $6, $7)`,
		input.IDCliente, input.IDEmpleado, input.Fecha, input.Ramos,
		nil, nil, nil,
	).Scan(&idVenta, &precioTotal, &mensaje)

	if err != nil {
		http.Error(w, "Error al crear venta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if mensaje != "Venta creada exitosamente" {
		http.Error(w, mensaje, http.StatusBadRequest)
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
