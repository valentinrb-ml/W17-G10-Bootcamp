package router

import (
	"github.com/go-chi/chi/v5"
	carryHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/carry"
	geographyHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/geography"
)

func MountGeographyRoutes(api chi.Router, hdGeography *geographyHandler.GeographyHandler, hdCarry *carryHandler.CarryHandler) {
	api.Route("/localities", func(r chi.Router) {
		r.Post("/", hdGeography.Create)
		r.Get("/reportSellers", hdGeography.CountSellersByLocality)
		r.Get("/reportCarries", hdCarry.ReportCarries)
	})
}
