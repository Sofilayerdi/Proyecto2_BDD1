package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"proy2-bck/db"
	"proy2-bck/handlers"
)

func main() {
	db.InitDB()
	db.InitDB_ORM()
	defer db.DB.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	r.Post("/login", handlers.Login)
	r.Post("/logout", handlers.Logout)

	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		r.Get("/productos", handlers.ListarProductos)
		r.Get("/productos/{id}", handlers.VerProducto)

		r.Group(func(r chi.Router) {
			r.Use(handlers.RolMiddleware("superadmin", "gerente", "vendedor"))
			r.Post("/productos", handlers.CrearProducto)
			r.Put("/productos/{id}", handlers.EditarProducto)
		})

		r.Group(func(r chi.Router) {
			r.Use(handlers.RolMiddleware("superadmin", "gerente"))
			r.Delete("/productos/{id}", handlers.EliminarProducto)
		})

		r.Get("/ramos", handlers.ListarRamos)
		r.Get("/ramos/{id}", handlers.VerRamo)

		r.Group(func(r chi.Router) {
			r.Use(handlers.RolMiddleware("superadmin", "gerente", "vendedor", "comprador"))
			r.Post("/ramos", handlers.CrearRamo)
		})

		r.Group(func(r chi.Router) {
			r.Use(handlers.RolMiddleware("superadmin", "gerente"))
			r.Delete("/ramos/{id}", handlers.EliminarRamo)
		})

		r.Get("/ventas", handlers.ListarVentas)
		r.Get("/ventas/{id}", handlers.VerVenta)

		r.Group(func(r chi.Router) {
			r.Use(handlers.RolMiddleware("superadmin", "gerente", "vendedor", "comprador"))
			r.Post("/ventas", handlers.CrearVenta)
		})

		r.Group(func(r chi.Router) {
			r.Use(handlers.RolMiddleware("superadmin", "gerente"))
			r.Delete("/ventas/{id}", handlers.EliminarVenta)
		})

		r.Get("/clientes", handlers.ListarClientes)
		r.Get("/empleados", handlers.ListarEmpleados)
		r.Get("/proveedores", handlers.ListarProveedores)

		r.Group(func(r chi.Router) {
			r.Use(handlers.RolMiddleware("superadmin", "gerente", "auditor"))
			r.Get("/reportes/ventas-mensuales", handlers.VentasMensuales)
			r.Get("/reportes/top-productos", handlers.TopProductosVendidos)
		})
	})

	log.Println("Servidor corriendo en http://0.0.0.0:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
