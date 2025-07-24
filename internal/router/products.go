package router

import (
	"github.com/go-chi/chi/v5"
	productHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product"
	productRecordHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_record"
)

func MountProductRoutes(api chi.Router, productHandler *productHandler.ProductHandler, productRecordHandler *productRecordHandler.ProductRecordHandler) {
	api.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.GetAll)
		r.Post("/", productHandler.Create)

		// Product records reporting path
		r.Get("/reportRecords", productRecordHandler.GetRecordsReport) // GET /api/v1/products/reportRecords?id=1

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", productHandler.GetByID)
			r.Patch("/", productHandler.Patch)
			r.Delete("/", productHandler.Delete)
		})
	})
}
