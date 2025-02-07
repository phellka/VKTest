package config

import (
	"os"
	"time"
)

const (
	HTTPTimeout    = 5 * time.Second
	UpdateInterval = 200 * time.Second
	UpdateWait     = 30 * time.Second
	PingInterval   = 15 * time.Second
)

var ServerURL = os.Getenv("BACKEND_URL")
var RabbitMQURL = os.Getenv("RabbitMQURL")

var PingLogsQueueName = os.Getenv("PingLogsQueueName")
var ContainersRequestQueueName = os.Getenv("ContainersRequestQueueName")
