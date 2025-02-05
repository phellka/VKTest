package service

import (
	"backend/internal/app/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var ErrContainerNotFound = errors.New("container not found")
var ErrPingLogNotFound = errors.New("pingLog not found")
var ErrPingLogNotSuccessful = errors.New("pingLog is not successful")
var ErrPingLogDNBContainer = errors.New("The provided PingLog does not belong to the specified container")
var ErrFailedUpdContainer = errors.New("failed to update container")
var ErrFailedCrtPinglog = errors.New("failed to create PingLog")

type Service struct {
	db *gorm.DB
}

func New(database *gorm.DB) *Service {
	return &Service{db: database}
}

func (s *Service) GetContainers() ([]models.Container, error) {
	var containers []models.Container
	result := s.db.Find(&containers)
	if result.Error != nil {
		return nil, result.Error
	}
	return containers, nil
}

func (s *Service) GetContainerByID(id int) (models.Container, error) {
	var container models.Container
	result := s.db.First(&container, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.Container{}, ErrContainerNotFound
		}
		return models.Container{}, result.Error
	}
	return container, nil
}

func (s *Service) GetContainerLastSuccessfulPing(id int) (models.PingLog, error) {
	var container models.Container
	if err := s.db.First(&container, id).Error; err != nil {
		return models.PingLog{}, ErrContainerNotFound
	}
	var pingLog models.PingLog
	result := s.db.Where("container_id = ? AND success = true", id).
		Order("timestamp DESC").
		First(&pingLog)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.PingLog{}, ErrPingLogNotFound
		}
		return models.PingLog{}, result.Error
	}
	return pingLog, nil
}

func (s *Service) PatchContainerLastSuccessfulPing(req models.UpdateContainerRequest) (models.Container, error) {
	var container models.Container
	if err := s.db.First(&container, req.ContainerID).Error; err != nil {
		return container, ErrContainerNotFound
	}

	var pingLog models.PingLog
	if err := s.db.First(&pingLog, req.LastSuccessfulPingId).Error; err != nil {
		return container, ErrPingLogNotFound
	}
	if pingLog.ContainerId != container.ID {
		return container, ErrPingLogDNBContainer
	}
	if !pingLog.Success {
		return container, ErrPingLogNotSuccessful
	}

	container.LastSuccessfulPingId = uint(req.LastSuccessfulPingId)
	if err := s.db.Save(&container).Error; err != nil {
		return container, ErrFailedUpdContainer
	}
	return container, nil
}

func (s *Service) CreatePingLog(pingLog models.PingLog) error {
	if pingLog.Timestamp.IsZero() {
		pingLog.Timestamp = time.Now()
	}
	var container models.Container
	if err := s.db.First(&container, pingLog.ContainerId).Error; err != nil {
		return ErrContainerNotFound
	}
	fmt.Println(container)
	if err := s.db.Create(&pingLog).Error; err != nil {
		return ErrFailedCrtPinglog
	}
	return nil
}
