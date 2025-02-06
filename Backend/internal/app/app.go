package app

import (
	"backend/internal/endpoint"
	"backend/internal/mw"
	"backend/internal/service"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	a.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	a.s = service.New(a.db)
	a.e = endpoint.New(a.s)

	a.echo = echo.New()
	a.echo.Use(mw.Database(a.db))

	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	a.setupRoutes()

	return a, nil
}

func (a *App) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on port %s\n", port)

	if err := a.echo.Start(":" + port); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (a *App) setupRoutes() {
	a.echo.GET("/containers", a.e.GetContainers)
	a.echo.GET("/containers/with-last-ping", a.e.GetContainersWithLastPing)
	a.echo.GET("/container", a.e.GetContainer)
	a.echo.GET("/container/lastsuccessful", a.e.GetContainerLastSuccessfulPing)
	a.echo.POST("/pinglog", a.e.PostPingLog)
}
