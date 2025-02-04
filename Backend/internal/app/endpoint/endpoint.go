package endpoint

import (
	"backend/internal/app/models"
	"backend/internal/app/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	s *service.Service
}

func New(s *service.Service) *Endpoint {
	return &Endpoint{s: s}
}

func (e *Endpoint) GetContainers(c echo.Context) error {
	containers, err := e.s.GetContainers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Code: http.StatusInternalServerError, Message: "Database error " + err.Error()})
	}
	return c.JSON(http.StatusOK, containers)
}

func (e *Endpoint) GetContainer(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid ID " + err.Error()})
	}
	container, err := e.s.GetContainerByID(id)
	if err != nil {
		if err == service.ErrContainerNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
		} else {
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
		}
	}
	return c.JSON(http.StatusOK, container)
}

func (e *Endpoint) GetContainerLastSuccessfulPing(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid ID " + err.Error()})
	}
	pingLog, err := e.s.GetContainerLastSuccessfulPing(id)
	if err != nil {
		switch err {
		case service.ErrContainerNotFound, service.ErrPingLogNotFound:
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
		}
	}
	return c.JSON(http.StatusOK, pingLog)
}

func (e *Endpoint) PatchContainerLastSuccessfulPing(c echo.Context) error {
	var req models.UpdateContainerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid request body " + err.Error()})
	}
	container, err := e.s.PatchContainerLastSuccessfulPing(req)
	if err != nil {
		switch err {
		case service.ErrContainerNotFound, service.ErrPingLogNotFound:
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
		case service.ErrPingLogNotSuccessful, service.ErrPingLogDNBContainer:
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
		}
	}
	return c.JSON(http.StatusOK, container)
}

func (e *Endpoint) PostPingLog(c echo.Context) error {
	var pingLog models.PingLog
	if err := c.Bind(&pingLog); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid request body " + err.Error()})
	}
	if err := e.s.CreatePingLog(pingLog); err != nil {
		if err == service.ErrContainerNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
		} else {
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
		}
	}
	return c.JSON(http.StatusCreated, pingLog)
}
