package server

import (
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

func (s *ServerChi) Run() error {
	// - loader

	buyerLd := loader.NewBuyerJSONFile("docs/db/buyers.json")
	buyerDb, err := buyerLd.Load()
	if err != nil {
		return
	}

	loadWarehouse := loader.NewWarehouseJSONFile("docs/db/warehouse.json")
	dbWarehouse, err := loadWarehouse.Load()
	if err != nil {
		return
	}

	sellerLd := loader.NewSellerJSONFile("docs/db/seller.json")
	sellerDb, err := sellerLd.Load()
	if err != nil {
		return
	}
  l := loader.NewEmployeeJSONFile("docs/db/employees.json")
	db, err := l.Load()
	if err != nil {
		return err
	}
	// - repository
	buyerRp := repository.NewBuyerRepository(buyerDb)
	// - service
	buyerSv := service.NewBuyerService(buyerRp)
	// - handler
	buyerHd := handler.NewBuyerHandler(buyerSv)

	rpWarehouse := repository.NewWarehouseMap(dbWarehouse)
	// - service
	svWarehouse := service.NewWarehouseDefault(rpWarehouse)
	// - handler
	hdWarehouse := handler.NewWarehouseDefault(svWarehouse)

	sellerRp := repository.NewSellerRepository(sellerDb)
	// - service
	sellerSv := service.NewSellerService(sellerRp)
	// - handler
	sellerHd := handler.NewSellerHandler(sellerSv)
  
  // - repository
  repo := repository.NewEmployeeMap()
	for _, emp := range db {
		_, _ = repo.Create(emp)
	}
  // - service
  svc := service.NewEmployeeDefault(repo)
  // - handler
  hd := handler.NewEmployeeHandler(svc)
  

	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)

	// - endpoints
	rt.Route("/api/v1/buyers", func(rt chi.Router) {
		rt.Post("/", buyerHd.Create())
		rt.Patch("/{id}", buyerHd.Update())
		rt.Delete("/{id}", buyerHd.Delete())
		rt.Get("/", buyerHd.FindAll())
		rt.Get("/{id}", buyerHd.FindById())
	})

	rt.Route("/warehouses", func(rt chi.Router) {
		rt.Post("/", hdWarehouse.Create)
		rt.Get("/", hdWarehouse.FindAll)
		rt.Get("/{id}", hdWarehouse.FindById)
		rt.Patch("/{id}", hdWarehouse.Update)
		rt.Delete("/{id}", hdWarehouse.Delete)
	})

	rt.Route("/seller", func(rt chi.Router) {
		rt.Post("/", sellerHd.Create())
		rt.Patch("/{id}", sellerHd.Update())
		rt.Delete("/{id}", sellerHd.Delete())
		rt.Get("/", sellerHd.FindAll())
		rt.Get("/{id}", sellerHd.FindById())
	})
  rt.Route("/api/v1/employees", func(rt chi.Router) {
		rt.Post("/", hd.Create)
		rt.Get("/", hd.GetAll)
		rt.Get("/{id}", hd.GetByID)
		rt.Patch("/{id}", hd.Update)
		rt.Delete("/{id}", hd.Delete)
	})

	fmt.Printf("Server running at http://localhost%s\n", a.serverAddress)

	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}
