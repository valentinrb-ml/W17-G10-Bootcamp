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
	ldBuyer := loader.NewBuyerJSONFile("docs/db/buyers.json")
	dbBuyer, err := ldBuyer.Load()
	if err != nil {
		return err
	}
	ldWarehouse := loader.NewWarehouseJSONFile("docs/db/warehouse.json")
	dbWarehouse, err := ldWarehouse.Load()
	if err != nil {
		return err
	}
	ldSeller := loader.NewSellerJSONFile("docs/db/seller.json")
	dbSeller, err := ldSeller.Load()
	if err != nil {
		return err
	}
	ldSection := loader.NewSectionJSONFile("docs/db/section.json")
	dbSection, err := ldSection.Load()
	if err != nil {
		return err
	}
	// - repository
	repoSection := repository.NewSectionMap(dbSection)
	repoSeller := repository.NewSellerRepository(dbSeller)
	repoBuyer := repository.NewBuyerRepository(dbBuyer)
	repoWarehouse := repository.NewWarehouseMap(dbWarehouse)
	repoProduct := repository.NewProductRepository(dbProduct)
	repoEmployee := repository.NewEmployeeRepository(mysql)

	// - service
	svcSeller := service.NewSellerService(repoSeller)
	svcBuyer := service.NewBuyerService(repoBuyer)
	svcSection := service.NewSectionServer(repoSection, repoWarehouse)
	svcProduct := service.NewProductService(repoProduct)
	svcWarehouse := service.NewWarehouseDefault(repoWarehouse)
	svcEmployee := service.NewEmployeeDefault(repoEmployee, repoWarehouse)

	// - handler
	hdBuyer := handler.NewBuyerHandler(svcBuyer)
	hdSection := handler.NewSectionHandler(svcSection)
	hdSeller := handler.NewSellerHandler(svcSeller)
	hdWarehouse := handler.NewWarehouseDefault(svcWarehouse)
	hdProduct := handler.NewProductHandler(svcProduct)
	hdEmployee := handler.NewEmployeeHandler(svcEmployee)

	// router
	rt := router.NewAPIRouter(hdBuyer, hdSection, hdSeller, hdWarehouse, hdProduct, hdEmployee)

	fmt.Printf("Server running at http://localhost%s\n", s.serverAddress)

	// run server
	err = http.ListenAndServe(s.serverAddress, rt)
	return err
}
