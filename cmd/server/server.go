package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/loader"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository"
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
func (s *ServerChi) Run() (err error) {
	ctx := context.Background()

	// dependencies

	// - loader
	ldProduct, err := loader.NewJSONFileProductLoader("docs/db/products.json")
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
	ldEmployee := loader.NewEmployeeJSONFile("docs/db/employees.json")
	dbEmployee, err := ldEmployee.Load()
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
	repoProduct, err := repository.NewProductRepository(ctx, ldProduct)
	if err != nil {
		return err
	}
	repoEmployee := repository.NewEmployeeMap(dbEmployee)

	// - service
	svcSeller := service.NewSellerService(repoSeller)
	svcBuyer := service.NewBuyerService(repoBuyer)
	svcSection := service.NewSectionServer(repoSection)
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
	rt := chi.NewRouter()
	rtProduct := hdProduct.Routes()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)

	rt.Mount("/api/v1/products", rtProduct)

	rt.Route("/api/v1/sections", func(r chi.Router) {
		r.Get("/", hdSection.FindAllSections())
		r.Get("/{id}", hdSection.FindById())
		r.Post("/", hdSection.CreateSection())
		r.Patch("/{id}", hdSection.UpdateSection())
		r.Delete("/{id}", hdSection.DeleteSection())

	})

	rt.Route("/api/v1/buyers", func(rt chi.Router) {
		rt.Post("/", hdBuyer.Create())
		rt.Patch("/{id}", hdBuyer.Update())
		rt.Delete("/{id}", hdBuyer.Delete())
		rt.Get("/", hdBuyer.FindAll())
		rt.Get("/{id}", hdBuyer.FindById())
	})

	rt.Route("/api/v1/warehouses", func(rt chi.Router) {
		rt.Post("/", hdWarehouse.Create)
		rt.Get("/", hdWarehouse.FindAll)
		rt.Get("/{id}", hdWarehouse.FindById)
		rt.Patch("/{id}", hdWarehouse.Update)
		rt.Delete("/{id}", hdWarehouse.Delete)
	})

	rt.Route("/api/v1/sellers", func(rt chi.Router) {
		rt.Post("/", hdSeller.Create())
		rt.Patch("/{id}", hdSeller.Update())
		rt.Delete("/{id}", hdSeller.Delete())
		rt.Get("/", hdSeller.FindAll())
		rt.Get("/{id}", hdSeller.FindById())
	})
	rt.Route("/api/v1/employees", func(rt chi.Router) {
		rt.Post("/", hdEmployee.Create)
		rt.Get("/", hdEmployee.GetAll)
		rt.Get("/{id}", hdEmployee.GetByID)
		rt.Patch("/{id}", hdEmployee.Update)
		rt.Delete("/{id}", hdEmployee.Delete)
	})

	fmt.Printf("Server running at http://localhost%s\n", s.serverAddress)

	// run server
	err = http.ListenAndServe(s.serverAddress, rt)
	return err
}
