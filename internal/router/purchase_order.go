package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
)

func MountPurchaseOrderRoutes(r chi.Router, h *handler.PurchaseOrderHandler) {
	// Ruta para creación de Purchase Orders
	r.Post("/purchaseOrders", h.Create)

	// Ruta para el reporte
	r.Get("/buyers/reportPurchaseOrders", h.GetReport)
}
