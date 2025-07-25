package server

import (
	"database/sql"
	"fmt"
	"net/http"

	secRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/section"

	productHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product"
	productRecordHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/product_record"
	productRepository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product"
	productRecordRepository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/product_record"
	productService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product"
	productRecordService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/product_record"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
	empHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/employee"
	inbHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/inbound_order"
	sellerHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/seller"
	warehouseHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	empRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"

	inbRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/inbound_order"
	sellerRepository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/seller"
	wRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/warehouse"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/router"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"

	buyerHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/buyer"
	buyerRepository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/buyer"
	buyerService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/buyer"
	empService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/employee"
	inbService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/inbound_order"
	sellerService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/seller"
	wService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/warehouse"

	geographyHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/geography"
	geographyRepository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
	geographyService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/geography"
)

type ConfigServerChi struct {
	ServerAddress string
}

func NewServerChi(cfg *ConfigServerChi) *ServerChi {
	defaultConfig := &ConfigServerChi{
		ServerAddress: ":8080",
	}
	if cfg != nil && cfg.ServerAddress != "" {
		defaultConfig.ServerAddress = cfg.ServerAddress
	}
	return &ServerChi{
		serverAddress: defaultConfig.ServerAddress,
	}
}

type ServerChi struct {
	serverAddress string
}

// Run is a method that runs the server
func (s *ServerChi) Run(mysql *sql.DB) (err error) {
	// - repository

	repoSection := secRepo.NewSectionRepository(mysql)
	repoSeller := sellerRepository.NewSellerRepository(mysql)
	repoBuyer := buyerRepository.NewBuyerRepository(mysql)
	repoWarehouse := wRepo.NewWarehouseRepository(mysql)
	repoProduct, err := productRepository.NewProductRepository(mysql)
	if err != nil {
		return err
	}
	repoEmployee := empRepo.NewEmployeeRepository(mysql)
	repoProductBatches := repository.NewProductBatchesRepository(mysql)
	repoCarry := repository.NewCarryRepository(mysql)
	repoGeography := geographyRepository.NewGeographyRepository(mysql)
	repoInboundOrder := inbRepo.NewInboundOrderRepository(mysql)
	repoPurchaseOrder := repository.NewPurchaseOrderRepository(mysql)
	repoProductRecord, err := productRecordRepository.NewProductRecordRepository(mysql)
	if err != nil {
		return err
	}

	// - service
	svcSeller := sellerService.NewSellerService(repoSeller, repoGeography)
	svcBuyer := buyerService.NewBuyerService(repoBuyer)
	svcSection := service.NewSectionServer(repoSection)
	svcProduct := productService.NewProductService(repoProduct)
	svcEmployee := empService.NewEmployeeDefault(repoEmployee, repoWarehouse)
	svcWarehouse := wService.NewWarehouseService(repoWarehouse)
	svcProductBatches := service.NewProductBatchesService(repoProductBatches)
	svcCarry := service.NewCarryService(repoCarry, repoGeography)
	svcGeography := geographyService.NewGeographyService(repoGeography)
	svcInboundOrder := inbService.NewInboundOrderService(repoInboundOrder, repoEmployee, repoWarehouse)
	svcPurchaseOrder := service.NewPurchaseOrderService(repoPurchaseOrder)
	svcProductRecord := productRecordService.NewProductRecordService(repoProductRecord)

	// - handler
	hdBuyer := buyerHandler.NewBuyerHandler(svcBuyer)
	hdSection := handler.NewSectionHandler(svcSection)
	hdSeller := sellerHandler.NewSellerHandler(svcSeller)
	hdCarry := handler.NewCarryHandler(svcCarry)
	hdWarehouse := warehouseHandler.NewWarehouseHandler(svcWarehouse)
	hdEmployee := empHandler.NewEmployeeHandler(svcEmployee)
	hdProduct := productHandler.NewProductHandler(svcProduct)
	hdProductBatches := handler.NewProductBatchesHandler(svcProductBatches)
	hdGeography := geographyHandler.NewGeographyHandler(svcGeography)
	hdInboundOrder := inbHandler.NewInboundOrderHandler(svcInboundOrder)
	hdPurchaseOrder := handler.NewPurchaseOrderHandler(svcPurchaseOrder)
	hdProductRecord := productRecordHandler.NewProductRecordHandler(svcProductRecord)

	// router
	rt := router.NewAPIRouter(
		hdBuyer, hdSection, hdSeller, hdWarehouse, hdEmployee,
		hdProduct, hdProductBatches, hdPurchaseOrder,
		hdGeography, hdInboundOrder, hdCarry, hdProductRecord,
	)

	fmt.Printf("Server running at http://localhost%s\n", s.serverAddress)

	// run server
	err = http.ListenAndServe(s.serverAddress, rt)
	return err
}
