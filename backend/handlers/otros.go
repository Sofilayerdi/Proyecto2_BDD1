package handlers

import (
	"encoding/json"
	"net/http"

	db "proy2-bck/db"
)

type Cliente struct {
	IDCliente int    `json:"id_cliente"`
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo"`
	Direccion string `json:"direccion"`
}

type Empleado struct {
	IDEmpleado int    `json:"id_empleado"`
	Nombre     string `json:"nombre"`
	Rol        string `json:"rol"`
}

type Proveedor struct {
	IDProveedor int    `json:"id_proveedor"`
	Nombre      string `json:"nombre"`
	Correo      string `json:"correo"`
	Telefono    string `json:"telefono"`
}

func ListarClientes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id_cliente, nombre, telefono, correo, direccion FROM cliente ORDER BY nombre`)
	if err != nil {
		http.Error(w, "Error al obtener clientes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []Cliente
	for rows.Next() {
		var c Cliente
		rows.Scan(&c.IDCliente, &c.Nombre, &c.Telefono, &c.Correo, &c.Direccion)
		lista = append(lista, c)
	}
	if lista == nil {
		lista = []Cliente{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func ListarEmpleados(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id_empleado, nombre, rol FROM empleado ORDER BY nombre`)
	if err != nil {
		http.Error(w, "Error al obtener empleados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []Empleado
	for rows.Next() {
		var e Empleado
		rows.Scan(&e.IDEmpleado, &e.Nombre, &e.Rol)
		lista = append(lista, e)
	}
	if lista == nil {
		lista = []Empleado{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func ListarProveedores(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id_proveedor, nombre, correo, telefono FROM proveedor ORDER BY nombre`)
	if err != nil {
		http.Error(w, "Error al obtener proveedores", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []Proveedor
	for rows.Next() {
		var p Proveedor
		rows.Scan(&p.IDProveedor, &p.Nombre, &p.Correo, &p.Telefono)
		lista = append(lista, p)
	}
	if lista == nil {
		lista = []Proveedor{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}
