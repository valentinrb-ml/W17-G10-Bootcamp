package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
)

// Monta las rutas relacionadas con la entidad Inbound_Orders en el router API principal.
func MountInboundOrderRoutes(api chi.Router, hd *handler.InboundOrderHandler) {
	api.Route("/inboundOrders", func(r chi.Router) {
		r.Post("/", hd.Create)
	})
	api.Get("/employees/reportInboundOrders", hd.Report)
}
