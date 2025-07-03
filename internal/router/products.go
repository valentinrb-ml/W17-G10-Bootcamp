package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
)

func MountProductRoutes(api chi.Router, hd *handler.ProductHandler) {
	api.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAll)
		r.Get("/{id}", hd.GetByID)
		r.Post("/", hd.Create)
		r.Patch("/{id}", hd.Patch)
		r.Delete("/{id}", hd.Delete)
	})
}