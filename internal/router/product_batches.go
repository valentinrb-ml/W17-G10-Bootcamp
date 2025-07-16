package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
)

func MountProductBatchesRoutes(api chi.Router, hd *handler.ProductBatchesHandler) {
	api.Route("/productBatches", func(r chi.Router) {
		r.Post("/", hd.CreateProductBatches)
	})
}
