package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
)

func MountSectionRoutes(api chi.Router, hd *handler.SectionDefault) {
	api.Route("/sections", func(r chi.Router) {
		r.Get("/", hd.FindAllSections)
		r.Get("/{id}", hd.FindById)
		r.Post("/", hd.CreateSection)
		r.Patch("/{id}", hd.UpdateSection)
		r.Delete("/{id}", hd.DeleteSection)
	})
}