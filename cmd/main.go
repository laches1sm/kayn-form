package main

import (
	"kayn-form/cmd/adapters"
	"kayn-form/cmd/httpserver"
	"log"
)

func main() {
	log := log.New(log.Writer(), "help-pix-", 0)
	adapter := adapters.NewKaynFormAdapter(log)

	server := httpserver.NewParrotServer(log, adapter)
	log.Print(`Creating server`)
	server.SetupRoutes()
	if err := server.Start(httpserver.ServerPort); err != nil {
		log.Println(err.Error())
	}
}
