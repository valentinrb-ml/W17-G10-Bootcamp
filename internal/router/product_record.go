package router

import (
	"github.com/go-chi/chi/v5"
	productRecordHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_record"
)

func MountProductRecordRoutes(api chi.Router, hd *productRecordHandler.ProductRecordHandler) {
	api.Route("/productRecords", func(r chi.Router) {
		r.Post("/", hd.Create)
	})
}
