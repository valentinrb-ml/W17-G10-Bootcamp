package router

import (
	"github.com/go-chi/chi/v5"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/employee"
)

// Monta las rutas relacionadas con la entidad Employee en el router API principal.
func MountEmployeeRoutes(api chi.Router, hd *handler.EmployeeHandler) {
	api.Route("/employees", func(r chi.Router) {
		r.Get("/", hd.GetAll)
		r.Get("/{id}", hd.GetByID)
		r.Post("/", hd.Create)
		r.Patch("/{id}", hd.Update)
		r.Delete("/{id}", hd.Delete)
	})
}
