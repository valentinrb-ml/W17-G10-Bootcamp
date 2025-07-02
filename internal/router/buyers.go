package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
)

func MountBuyerRoutes(api chi.Router, hd *handler.BuyerHandler) {
	api.Route("/buyers", func(r chi.Router) {
		r.Get("/", hd.FindAll)
		r.Get("/{id}", hd.FindById)
		r.Post("/", hd.Create)
		r.Patch("/{id}", hd.Update)
		r.Delete("/{id}", hd.Delete)
	})
}