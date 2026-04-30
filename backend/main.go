package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	db "tu-modulo/db"
	"tu-modulo/handlers"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/productos", handlers.ListarProductos)
	r.Get("/productos/{id}", handlers.VerProducto)
	r.Post("/productos", handlers.CrearProducto)
	r.Put("/productos/{id}", handlers.EditarProducto)
	r.Delete("/productos/{id}", handlers.EliminarProducto)

	r.Post("/ventas", handlers.CrearVenta)
	r.Get("/ventas", handlers.ListarVentas)
	r.PUT("/ventas/{id}", handlers.EditarVenta)
	r.DELETE("/ventas/{id}", handlers.EliminarVenta)

	r.Get("/reportes/ventas-por-empleado", handlers.VentasPorEmpleado)
	r.Get("/reportes/productos-en-ramos", handlers.ProductosEnRamos)
	r.Get("/reportes/inventario", handlers.Inventario)
	r.Get("/reportes/vista-ventas", handlers.VistaVentas)

	log.Println("Servidor corriendo en http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
