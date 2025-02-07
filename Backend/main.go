package main

import (
	"backend/internal/app"
	"log"
	"os"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("Error initializing app: %v", err)
		os.Exit(1)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
		os.Exit(1)
	}

	log.Println("Server started successfully on port 8080")
}
