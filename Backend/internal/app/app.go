package app

import (
	"backend/internal/app/endpoint"
	"backend/internal/app/mw"
	"backend/internal/app/service"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	e    *endpoint.Endpoint
	s    *service.Service
	echo *echo.Echo
	db   *gorm.DB
}

func New() (*App, error) {
	a := &App{}

	dsn := "host=localhost user=postgres password=postgres dbname=container_monitoring port=5435"
	var err error
	a.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	a.s = service.New(a.db)
	a.e = endpoint.New(a.s)

	a.echo = echo.New()
	a.echo.Use(mw.Database(a.db))

	a.setupRoutes()

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("Server running")
	if err := a.echo.Start(":8080"); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (a *App) setupRoutes() {
	a.echo.GET("/containers", a.e.GetContainers)
	a.echo.GET("/container", a.e.GetContainer)
	a.echo.GET("/container/lastsuccessful", a.e.GetContainerLastSuccessfulPing)
	a.echo.POST("/pinglog", a.e.PostPingLog)
	a.echo.PATCH("/container", a.e.PatchContainerLastSuccessfulPing)
}
