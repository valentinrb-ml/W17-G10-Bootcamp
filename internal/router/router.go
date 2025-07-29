package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
	buyerHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/buyer"
	carryHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/carry"
	empHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/employee"
	geographyHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/geography"
	inbHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/inbound_order"
	productHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product"
	ProductRecordHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_record"
	purchaseOrderHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/purchase_order"
	sectionHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/section"
	sellerHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/seller"
	warehouseHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
)

func NewAPIRouter(
	hdBuyer *buyerHandler.BuyerHandler,
	hdSection *sectionHandler.SectionDefault,
	hdSeller *sellerHandler.SellerHandler,
	hdWarehouse *warehouseHandler.WarehouseHandler,
	hdEmployee *empHandler.EmployeeHandler,
	hdProduct *productHandler.ProductHandler,
	hdProductBatches *handler.ProductBatchesHandler,
	hdPurchaseOrder *purchaseOrderHandler.PurchaseOrderHandler,
	hdGeography *geographyHandler.GeographyHandler,
	hdInboundOrder *inbHandler.InboundOrderHandler,
	hdCarry *carryHandler.CarryHandler,
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
