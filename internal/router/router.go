package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
	buyerHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/buyer"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
)

func NewAPIRouter(
	hdBuyer *buyerHandler.BuyerHandler,
	hdSection *handler.SectionDefault,
	hdSeller *handler.SellerHandler,
	hdWarehouse *handler.WarehouseHandler,
	hdProduct *handler.ProductHandler,
	hdEmployee *handler.EmployeeHandler,
	hdProductBatches *handler.ProductBatchesHandler,
	hdPurchaseOrder *handler.PurchaseOrderHandler,
	hdGeography *handler.GeographyHandler,
	hdInboundOrder *handler.InboundOrderHandler,
	hdCarry *handler.CarryHandler,
	hdProductRecord *handler.ProductRecordHandler,
) *chi.Mux {
	root := chi.NewRouter()
	root.Use(middleware.Logger, middleware.Recoverer)

	root.MethodNotAllowed(httputil.MethodNotAllowedHandler)
	root.NotFound(httputil.NotFoundHandler)

	root.Route("/api/v1", func(api chi.Router) {
		MountProductRoutes(api, hdProduct, hdProductRecord)
		MountSectionRoutes(api, hdSection, hdProductBatches)
		MountBuyerRoutes(api, hdBuyer)
		MountWarehouseRoutes(api, hdWarehouse)
		MountSellerRoutes(api, hdSeller)
		MountEmployeeRoutes(api, hdEmployee)
		MountProductBatchesRoutes(api, hdProductBatches)
		MountPurchaseOrderRoutes(api, hdPurchaseOrder)
		MountCarryRoutes(api, hdCarry)
		MountGeographyRoutes(api, hdGeography, hdCarry)
		MountInboundOrderRoutes(api, hdInboundOrder)
		MountProductRecordRoutes(api, hdProductRecord)
	})

	return root
}
