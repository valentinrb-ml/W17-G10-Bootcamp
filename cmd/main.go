package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/cmd/server"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Print("Could not load .env file, continuing with system variables only")
    }
	
	// env
	mysql, err := database.InitMysqlDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer mysql.Close()

	// app
	// - config
	cfg := &server.ConfigServerChi{
		ServerAddress: ":8080",
	}
	app := server.NewServerChi(cfg)
	// - run
	if err := app.Run(mysql); err != nil {
		fmt.Println(err)
		return
	}
}
