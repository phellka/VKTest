package config

import (
	"os"
	"time"
)

const (
	HTTPTimeout    = 5 * time.Second
	UpdateInterval = 100 * time.Second
	PingInterval   = 10 * time.Second
)

var ServerURL = os.Getenv("BACKEND_URL")
