package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
)



func NewAPIRouter(
    hdBuyer *handler.BuyerHandler,
    hdSection *handler.SectionDefault,
    hdSeller *handler.SellerHandler,
    hdWarehouse *handler.WarehouseHandler,
    hdProduct *handler.ProductHandler,
    hdEmployee *handler.EmployeeHandler,
    ) *chi.Mux{
    root := chi.NewRouter()
	root.Use(middleware.Logger, middleware.Recoverer)

	root.Route("/api/v1", func(api chi.Router) {
        MountProductRoutes(api, hdProduct)
        MountSectionRoutes(api, hdSection)
        MountBuyerRoutes(api, hdBuyer)
        MountWarehouseRoutes(api, hdWarehouse)
        MountSellerRoutes(api, hdSeller)
        MountEmployeeRoutes(api, hdEmployee)
    })

	return root
}