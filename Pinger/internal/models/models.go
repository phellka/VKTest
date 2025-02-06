package models

import "time"

type Container struct {
	ID                   uint   `json:"id"`
	Ip                   string `json:"ip"`
	Name                 string `json:"name"`
	LastSuccessfulPingId uint   `json:"lastsuccessfulpingid"`
}

type PostPingLog struct {
	ContainerId uint      `json:"containerid"`
	Timestamp   time.Time `json:"timestamp"`
	Success     bool      `json:"success"`
}
