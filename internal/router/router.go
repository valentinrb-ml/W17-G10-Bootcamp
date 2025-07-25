package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
	buyerHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/buyer"
	empHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/employee"
	geographyHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/geography"
	inbHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/inbound_order"
	productHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product"
	ProductRecordHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_record"
	sellerHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/seller"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
)

func NewAPIRouter(
	hdBuyer *buyerHandler.BuyerHandler,
	hdSection *handler.SectionDefault,
	hdSeller *sellerHandler.SellerHandler,
	hdWarehouse *handler.WarehouseHandler,
	hdEmployee *empHandler.EmployeeHandler,
	hdProduct *productHandler.ProductHandler,
	hdProductBatches *handler.ProductBatchesHandler,
	hdPurchaseOrder *handler.PurchaseOrderHandler,
	hdGeography *geographyHandler.GeographyHandler,
	hdInboundOrder *inbHandler.InboundOrderHandler,
	hdCarry *handler.CarryHandler,
	hdProductRecord *ProductRecordHandler.ProductRecordHandler,
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
