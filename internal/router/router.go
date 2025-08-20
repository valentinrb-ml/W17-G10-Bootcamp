package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	buyerHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/buyer"
	carryHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/carry"
	empHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/employee"
	geographyHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/geography"
	inbHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/inbound_order"
	productHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product"
	productBatchHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_batch"
	ProductRecordHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_record"
	purchaseOrderHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/purchase_order"
	sectionHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/section"
	sellerHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/seller"
	warehouseHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
)

func NewAPIRouter(
	hdBuyer *buyerHandler.BuyerHandler,
	hdSection *sectionHandler.SectionDefault,
	hdSeller *sellerHandler.SellerHandler,
	hdWarehouse *warehouseHandler.WarehouseHandler,
	hdEmployee *empHandler.EmployeeHandler,
	hdProduct *productHandler.ProductHandler,
	hdProductBatches *productBatchHandler.ProductBatchesHandler,
	hdPurchaseOrder *purchaseOrderHandler.PurchaseOrderHandler,
	hdGeography *geographyHandler.GeographyHandler,
	hdInboundOrder *inbHandler.InboundOrderHandler,
	hdCarry *carryHandler.CarryHandler,
	hdProductRecord *ProductRecordHandler.ProductRecordHandler,
	appLogger logger.Logger,
) *chi.Mux {
	root := chi.NewRouter()

	// Basic middleware with logging
	root.Use(middleware.RequestID)
	root.Use(logger.LoggingMiddleware(appLogger))
	root.Use(middleware.Recoverer)

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
