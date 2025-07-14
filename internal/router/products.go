package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/middleware"
)

func MountProductRoutes(api chi.Router, hd *handler.ProductHandler) {
	api.Route("/products", func(r chi.Router) {
		r.With(middleware.OnlyGET).Get("/", hd.GetAll)
		r.With(middleware.OnlyPOST).Post("/", hd.Create)

		r.Route("/{id}", func(r chi.Router) {
			r.With(middleware.OnlyGET).Get("/", hd.GetByID)
			r.With(middleware.OnlyPATCH).Patch("/", hd.Patch)
			r.With(middleware.OnlyDELETE).Delete("/", hd.Delete)
		})
	})
}
