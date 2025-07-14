package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
)

func MountProductRoutes(api chi.Router, hd *handler.ProductHandler) {
	api.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAll)
		r.Post("/", hd.Create)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", hd.GetByID)
			r.Patch("/", hd.Patch)
			r.Delete("/", hd.Delete)
		})
	})
}
