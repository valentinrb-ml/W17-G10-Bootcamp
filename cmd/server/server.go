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

	repoSection := repository.NewSectionMap(mysql)
	repoSeller := repository.NewSellerRepository(mysql)
	repoBuyer := repository.NewBuyerRepository(mysql)
	repoWarehouse := repository.NewWarehouseRepository(mysql)
	repoProduct := repository.NewProductRepository(dbProduct)
	repoEmployee := repository.NewEmployeeRepository(mysql)
	repoInboundOrder := repository.NewInboundOrderRepository(mysql)

	// - service
	svcSeller := service.NewSellerService(repoSeller)
	svcBuyer := service.NewBuyerService(repoBuyer)
	svcSection := service.NewSectionServer(repoSection)
	svcProduct := service.NewProductService(repoProduct)
	// svcWarehouse := service.NewWarehouseDefault(repoWarehouse)
	svcWarehouse := service.NewWarehouseService(repoWarehouse)
	svcEmployee := service.NewEmployeeDefault(repoEmployee, repoWarehouse)
	svcInboundOrder := service.NewInboundOrderService(repoInboundOrder, repoEmployee, repoWarehouse)

	// - handler
	hdBuyer := handler.NewBuyerHandler(svcBuyer)
	hdSection := handler.NewSectionHandler(svcSection)
	hdSeller := handler.NewSellerHandler(svcSeller)
	// hdWarehouse := handler.NewWarehouseDefault(svcWarehouse)
	hdWarehouse := handler.NewWarehouseHandler(svcWarehouse)
	hdProduct := handler.NewProductHandler(svcProduct)
	hdEmployee := handler.NewEmployeeHandler(svcEmployee)
	hdInboundOrder := handler.NewInboundOrderHandler(svcInboundOrder)

	// router
	rt := router.NewAPIRouter(hdBuyer, hdSection, hdSeller, hdWarehouse, hdProduct, hdEmployee, hdInboundOrder)

	fmt.Printf("Server running at http://localhost%s\n", s.serverAddress)

	// run server
	err = http.ListenAndServe(s.serverAddress, rt)
	return err
}
