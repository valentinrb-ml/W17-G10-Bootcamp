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
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"

)

type ConfigServerChi struct {
	ServerAddress string
}

func NewServerChi(cfg *ConfigServerChi) *ServerChi {
	defaultConfig := &ConfigServerChi{
		ServerAddress: ":8080",
	}
	if cfg != nil {
		if cfg.ServerAddress != "" {
			defaultConfig.ServerAddress = cfg.ServerAddress
		}
	}

	return &ServerChi{
		serverAddress: defaultConfig.ServerAddress,
	}
}

type ServerChi struct {
	serverAddress string
}

// Run is a method that runs the server
func (a *ServerChi) Run() (err error) {
	// dependencies
	// - loader

	buyerLd := loader.NewBuyerJSONFile("docs/db/buyers.json")
	buyerDb, err := buyerLd.Load()


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

  fmt.Printf("Server running at http://localhost%s\n", a.serverAddress)
  
	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}
