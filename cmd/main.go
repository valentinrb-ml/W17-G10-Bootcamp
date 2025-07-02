package main

import (
	"fmt"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/cmd/server"
)

func main() {
	// env
	// ...

	// app
	// - config
	cfg := &server.ConfigServerChi{
		ServerAddress: ":8080",
	}
	app := server.NewServerChi(cfg)
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
