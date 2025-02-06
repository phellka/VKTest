package main

import (
	"log"
	"pinger/internal/app"
	"pinger/internal/service"
)

func main() {
	log.Println("Starting the application...")

	svc := service.NewService()

	app := app.NewApp(svc)
	app.Start()
}
