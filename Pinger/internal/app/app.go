package app

import (
	"log"
	"pinger/internal/config"
	"pinger/internal/service"
	"time"
)

type App struct {
	service *service.Service
}

func NewApp(service *service.Service) *App {
	return &App{
		service: service,
	}
}

func (a *App) Start() {
	log.Println("Starting pinger...")

	go a.updateContainersLoop()

	go a.pingLoop()

	select {}
}

func (a *App) updateContainersLoop() {
	log.Println("Starting container update loop")
	for {
		err := a.service.FetchContainers()
		if err != nil {
			log.Println("Error fetching containers:", err)
		}

		time.Sleep(config.UpdateInterval)
	}
}

func (a *App) pingLoop() {
	log.Println("Starting container ping loop")
	for {
		localContainers := a.service.GetContainers()

		for _, container := range localContainers {
			pingLog, err := a.service.PingContainer(container)
			if err != nil {
				log.Printf("Error pinging %s: %v", container.Ip, err)
			}

			if err := a.service.SendPingLog(pingLog); err != nil {
				log.Printf("Error sending ping log for %s: %v", container.Ip, err)
			}
		}

		time.Sleep(config.PingInterval)
	}
}
