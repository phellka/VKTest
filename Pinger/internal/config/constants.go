package config

import "time"

const (
	ServerURL      = "http://localhost:8080"
	HTTPTimeout    = 5 * time.Second
	UpdateInterval = 30 * time.Second
	PingInterval   = 5 * time.Second
)
