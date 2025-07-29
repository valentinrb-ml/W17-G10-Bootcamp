package router

import (
	"github.com/go-chi/chi/v5"
	productBatchHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_batch"
)

func MountProductBatchesRoutes(api chi.Router, hd *productBatchHandler.ProductBatchesHandler) {
	api.Route("/productBatches", func(r chi.Router) {
		r.Post("/", hd.CreateProductBatches)
	})
}
