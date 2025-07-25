package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/carry"
)

func MountCarryRoutes(api chi.Router, hd *handler.CarryHandler) {
	api.Route("/carries", func(r chi.Router) {
		r.Post("/", hd.Create)
	})
}