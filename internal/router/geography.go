package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
)

func MountGeographyRoutes(api chi.Router, hd *handler.GeographyHandler) {
	api.Route("/localities", func(r chi.Router) {
		r.Post("/", hd.Create)
	})
}
