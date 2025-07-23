package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
	empHandler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	empRepo "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/employee"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/router"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	empService "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/employee"
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
	repoSection := repository.NewSectionRepository(mysql)
	repoSeller := repository.NewSellerRepository(mysql)
	repoBuyer := repository.NewBuyerRepository(mysql)
	repoWarehouse := repository.NewWarehouseRepository(mysql)
	repoProduct, err := repository.NewProductRepository(mysql)
	if err != nil {
		return err
	}
	repoEmployee := empRepo.NewEmployeeRepository(mysql)
	repoProductBatches := repository.NewProductBatchesRepository(mysql)
	repoCarry := repository.NewCarryRepository(mysql)
	repoGeography := repository.NewGeographyRepository(mysql)
	repoInboundOrder := repository.NewInboundOrderRepository(mysql)
	repoPurchaseOrder := repository.NewPurchaseOrderRepository(mysql)
	repoProductRecord, err := repository.NewProductRecordRepository(mysql)
	if err != nil {
		return err
	}

	// - service
	svcSeller := service.NewSellerService(repoSeller, repoGeography)
	svcBuyer := service.NewBuyerService(repoBuyer)
	svcSection := service.NewSectionServer(repoSection)
	svcProduct := service.NewProductService(repoProduct)
	svcWarehouse := service.NewWarehouseService(repoWarehouse)
	svcEmployee := empService.NewEmployeeDefault(repoEmployee, repoWarehouse)
	svcProductBatches := service.NewProductBatchesService(repoProductBatches)
	svcCarry := service.NewCarryService(repoCarry, repoGeography)
	svcGeography := service.NewGeographyService(repoGeography)
	svcInboundOrder := service.NewInboundOrderService(repoInboundOrder, repoEmployee, repoWarehouse)
	svcPurchaseOrder := service.NewPurchaseOrderService(repoPurchaseOrder)
	svcProductRecord := service.NewProductRecordService(repoProductRecord)

	// - handler
	hdBuyer := handler.NewBuyerHandler(svcBuyer)
	hdSection := handler.NewSectionHandler(svcSection)
	hdSeller := handler.NewSellerHandler(svcSeller)
	hdCarry := handler.NewCarryHandler(svcCarry)
	hdWarehouse := handler.NewWarehouseHandler(svcWarehouse)
	hdProduct := handler.NewProductHandler(svcProduct)
	hdEmployee := empHandler.NewEmployeeHandler(svcEmployee)
	hdProductBatches := handler.NewProductBatchesHandler(svcProductBatches)
	hdGeography := handler.NewGeographyHandler(svcGeography)
	hdInboundOrder := handler.NewInboundOrderHandler(svcInboundOrder)
	hdPurchaseOrder := handler.NewPurchaseOrderHandler(svcPurchaseOrder)
	hdProductRecord := handler.NewProductRecordHandler(svcProductRecord)

	// router
	rt := router.NewAPIRouter(
		hdBuyer, hdSection, hdSeller, hdWarehouse, hdProduct,
		hdEmployee, hdProductBatches, hdPurchaseOrder,
		hdGeography, hdInboundOrder, hdCarry, hdProductRecord,
	)

	fmt.Printf("Server running at http://localhost%s\n", s.serverAddress)

	// run server
	err = http.ListenAndServe(s.serverAddress, rt)
	return err
}
