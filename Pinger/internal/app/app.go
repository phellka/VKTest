package app

import (
	"fmt"
	"log"
	"pinger/internal/config"
	"pinger/internal/service"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type App struct {
	service               *service.Service
	amqpConnect           *amqp.Connection
	amqpChannelPingLogs   *amqp.Channel
	amqpChannelContainers *amqp.Channel
	updateWG              sync.WaitGroup
	startOnce             sync.Once
}

func NewApp() (*App, error) {
	app := &App{}
	var err error
	app.amqpConnect, err = amqp.Dial(config.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	app.amqpChannelPingLogs, err = app.amqpConnect.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel to RabbitMQ: %w", err)
	}
	app.amqpChannelContainers, err = app.amqpConnect.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel to RabbitMQ: %w", err)
	}

	app.service = service.NewService(app.amqpChannelPingLogs, app.amqpChannelContainers)
	return app, nil
}

func (a *App) Start() {
	log.Println("Starting pinger...")
	defer a.Close()
	a.updateWG.Add(1)
	go a.updateContainersLoop()
	a.updateWG.Wait()
	go a.pingLoop()

	select {}
}

func (a *App) Close() {
	if a.amqpChannelPingLogs != nil {
		a.amqpChannelPingLogs.Close()
	}
	if a.amqpConnect != nil {
		a.amqpConnect.Close()
	}
}

func (a *App) updateContainersLoop() {
	log.Println("Starting container update loop")
	for {
		err := a.service.FetchContainers()
		if err != nil {
			log.Println("Error fetching containers:", err)
			continue
		} else {
			a.startOnce.Do(func() {
				log.Println("First successful container update")
				a.updateWG.Done()
			})
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

			if err := a.service.SendPingLogToQueue(pingLog); err != nil {
				log.Printf("Error sending ping log for %s: %v", container.Ip, err)
			}
		}

		time.Sleep(config.PingInterval)
	}
}
