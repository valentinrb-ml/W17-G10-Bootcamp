package server

import (
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
	// 1) Loader
	l := loader.NewEmployeeJSONFile("docs/db/employees.json")
	db, err := l.Load()
	if err != nil {
		return err
	}

	// 2) Repo en memoria
	repo := repository.NewEmployeeMap()
	for _, emp := range db {
		_, _ = repo.Create(emp)
	}

	// 3) Service
	svc := service.NewEmployeeDefault(repo)

	// 4) Handler
	hd := handler.NewEmployeeHandler(svc)

	// 5) Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/api/v1/employees", func(r chi.Router) {
		r.Post("/", hd.Create)
		r.Get("/", hd.GetAll)
		r.Get("/{id}", hd.GetByID)
		r.Patch("/{id}", hd.Update)
		r.Delete("/{id}", hd.Delete)
	})

	// 6) Run server
	return http.ListenAndServe(s.serverAddress, r)
}
