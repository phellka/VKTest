package models

import "time"

type Container struct {
	ID   uint   `gorm:"primaryKey"`
	Ip   string `gorm:"unique;not null;"`
	Name string
}

type ContainerWithPingTime struct {
	ID        uint
	Ip        string
	Name      string
	Timestamp *time.Time
	Pingtime  *float64
}

type PingLog struct {
	ID          uint      `gorm:"primaryKey"`
	ContainerId uint      `gorm:"not null"`
	Timestamp   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Success     bool
	Pingtime    *float64
}

type ErrorResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}
