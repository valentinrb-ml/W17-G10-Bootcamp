package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/loader"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/router"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
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
	ctx := context.Background()

	ldProduct, err := loader.NewJSONFileProductLoader("docs/db/products.json")
	if err != nil {
		return err
	}
	dbProduct, err := ldProduct.Load(ctx)
	if err != nil {
		return err
	}

	// - repository

	repoSection := repository.NewSectionRepository(mysql)
	repoSeller := repository.NewSellerRepository(mysql)
	repoBuyer := repository.NewBuyerRepository(mysql)
	repoWarehouse := repository.NewWarehouseRepository(mysql)
	repoProduct := repository.NewProductRepository(dbProduct)
	repoEmployee := repository.NewEmployeeRepository(mysql)
	repoProductBatches := repository.NewProductBatchesRepository(mysql)
	repoCarry := repository.NewCarryRepository(mysql)
	repoGeography := repository.NewGeographyRepository(mysql)
	repoInboundOrder := repository.NewInboundOrderRepository(mysql)
	repoPurchaseOrder := repository.NewPurchaseOrderRepository(mysql)


	// - service
	svcSeller := service.NewSellerService(repoSeller, repoGeography)
	svcBuyer := service.NewBuyerService(repoBuyer)
	svcSection := service.NewSectionServer(repoSection)
	svcProduct := service.NewProductService(repoProduct)
	svcWarehouse := service.NewWarehouseService(repoWarehouse)
	svcEmployee := service.NewEmployeeDefault(repoEmployee, repoWarehouse)
	svcProductBatches := service.NewProductBatchesService(repoProductBatches)
	svcCarry := service.NewCarryService(repoCarry)
	svcGeography := service.NewGeographyService(repoGeography)
	svcInboundOrder := service.NewInboundOrderService(repoInboundOrder, repoEmployee, repoWarehouse)
	svcPurchaseOrder := service.NewPurchaseOrderService(repoPurchaseOrder)


	// - handler
	hdBuyer := handler.NewBuyerHandler(svcBuyer)
	hdSection := handler.NewSectionHandler(svcSection)
	hdSeller := handler.NewSellerHandler(svcSeller)
	hdCarry := handler.NewCarryHandler(svcCarry)
	hdWarehouse := handler.NewWarehouseHandler(svcWarehouse)
	hdProduct := handler.NewProductHandler(svcProduct)
	hdEmployee := handler.NewEmployeeHandler(svcEmployee)
	hdProductBatches := handler.NewProductBatchesHandler(svcProductBatches)
	hdGeography := handler.NewGeographyHandler(svcGeography)
	hdInboundOrder := handler.NewInboundOrderHandler(svcInboundOrder)
	hdPurchaseOrder := handler.NewPurchaseOrderHandler(svcPurchaseOrder)

	// router
	rt := router.NewAPIRouter(hdBuyer, hdSection, hdSeller, hdWarehouse, hdProduct, hdEmployee, hdGeography, hdInboundOrder, hdCarry, hdPurchaseOrder)

	fmt.Printf("Server running at http://localhost%s\n", s.serverAddress)

	// run server
	err = http.ListenAndServe(s.serverAddress, rt)
	return err
}
