package models

import "time"

type Container struct {
	ID                   uint   `gorm:"primaryKey"`
	Ip                   string `gorm:"unique;not null;"`
	Name                 string
	LastSuccessfulPingId uint
}

type PingLog struct {
	ID          uint      `gorm:"primaryKey"`
	ContainerId uint      `gorm:"not null"`
	Timestamp   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Success     bool
}

type ErrorResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

type UpdateContainerRequest struct {
	ContainerID          uint `json:"containerid"`
	LastSuccessfulPingId uint `json:"lastsuccessfulpingid"`
}
