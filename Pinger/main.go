package main

import (
	"log"
	"pinger/internal/app"
)

func main() {
	log.Println("Starting the application...")

	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	app.Start()
}
