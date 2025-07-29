package router

import (
	"github.com/go-chi/chi/v5"
	productBatchHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_batch"
	sectionHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/section"
)

func MountSectionRoutes(api chi.Router, hd *sectionHandler.SectionDefault, hdPB *productBatchHandler.ProductBatchesHandler) {
	api.Route("/sections", func(r chi.Router) {
		r.Get("/", hd.FindAllSections)
		r.Post("/", hd.CreateSection)
		r.Get("/{id}", hd.FindById)
		r.Patch("/{id}", hd.UpdateSection)
		r.Delete("/{id}", hd.DeleteSection)

		//Product Batches report
		r.Get("/reportProduct", hdPB.GetReportProduct)
	})
}
