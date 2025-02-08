package service

import (
	"backend/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrContainerNotFound = errors.New("container not found")
var ErrPingLogNotFound = errors.New("pingLog not found")
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

func (s *Service) GetContainersWithLastPing() ([]models.ContainerWithPingTime, error) {
	var containers []models.ContainerWithPingTime
	result := s.db.Table("containers AS c").
		Select("c.*, pl.timestamp, pl.pingtime").
		Joins("LEFT JOIN ( " +
			"SELECT pl.container_id, pl.timestamp, pl.pingtime " +
			"FROM ping_logs pl " +
			"INNER JOIN ( " +
			"SELECT container_id, MAX(timestamp) AS max_timestamp " +
			"FROM ping_logs " +
			"WHERE success = true " +
			"GROUP BY container_id " +
			") AS max_pl ON pl.container_id = max_pl.container_id AND pl.timestamp = max_pl.max_timestamp " +
			") AS pl ON c.id = pl.container_id ").
		Scan(&containers)

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

func (s *Service) CreatePingLog(pingLog models.PingLog) error {
	if pingLog.Timestamp.IsZero() {
		pingLog.Timestamp = time.Now()
	}
	var container models.Container
	if err := s.db.First(&container, pingLog.ContainerId).Error; err != nil {
		return ErrContainerNotFound
	}
	if err := s.db.Create(&pingLog).Error; err != nil {
		return ErrFailedCrtPinglog
	}
	return nil
}
