package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pinger/internal/config"
	"pinger/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/streadway/amqp"
)

type Service struct {
	client                 *http.Client
	containers             []models.Container
	mutex                  sync.Mutex
	amqpChannelPingLogs    *amqp.Channel
	channelPingLogsMutex   sync.Mutex
	amqpChannelContainers  *amqp.Channel
	channelContainersMutex sync.Mutex
}

func NewService(amqpChannelPingLogs *amqp.Channel, amqpChannelContainers *amqp.Channel) *Service {
	return &Service{
		client:                &http.Client{Timeout: config.HTTPTimeout},
		amqpChannelPingLogs:   amqpChannelPingLogs,
		amqpChannelContainers: amqpChannelContainers,
	}
}

func (s *Service) FetchContainers() error {
	s.channelContainersMutex.Lock()
	defer s.channelContainersMutex.Unlock()
	corrID := uuid.New().String()

	replyQueue, err := s.amqpChannelPingLogs.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare reply queue: %w", err)
	}
	msgs, err := s.amqpChannelPingLogs.Consume(
		replyQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume from reply queue: %w", err)
	}

	requestQueue, err := s.amqpChannelPingLogs.QueueDeclare(
		config.ContainersRequestQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}
	err = s.amqpChannelPingLogs.Publish(
		"",
		requestQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          []byte(fmt.Sprintf(`{"request": "get_containers", "corrID": "%s"}`, corrID)),
			ReplyTo:       replyQueue.Name,
			CorrelationId: corrID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	timeout := time.After(config.UpdateWait)

	for {
		select {
		case msg := <-msgs:
			if msg.CorrelationId == corrID {
				var containers []models.Container
				if err := json.Unmarshal(msg.Body, &containers); err != nil {
					return fmt.Errorf("failed to unmarshal containers response: %w", err)
				}
				_, err := s.amqpChannelPingLogs.QueueDelete(replyQueue.Name, false, false, false)
				if err != nil {
					log.Printf("failed to delete reply queue: %v", err)
				}
				s.mutex.Lock()
				s.containers = containers
				s.mutex.Unlock()
				return nil
			}
		case <-timeout:
			_, err := s.amqpChannelPingLogs.QueueDelete(replyQueue.Name, false, false, false)
			if err != nil {
				log.Printf("failed to delete reply queue: %v", err)
			}
			return fmt.Errorf("timeout waiting for response")
		}
	}
}

func (s *Service) SendPingLogToQueue(pingLog models.PostPingLog) error {
	s.channelPingLogsMutex.Lock()
	defer s.channelPingLogsMutex.Unlock()
	q, err := s.amqpChannelPingLogs.QueueDeclare(
		config.PingLogsQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	jsonPingLog, err := json.Marshal(pingLog)
	if err != nil {
		return fmt.Errorf("failed to marshal ping log: %w", err)
	}

	err = s.amqpChannelPingLogs.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonPingLog,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (s *Service) PingContainer(container models.Container) (models.PostPingLog, error) {
	var pingLog models.PostPingLog
	pingLog.ContainerId = container.ID
	pingLog.Timestamp = time.Now().UTC()

	pinger, err := probing.NewPinger(container.Ip)
	if err != nil {
		return pingLog, fmt.Errorf("Error creating pinger for %s: %v", container.Ip, err)
	}

	pinger.Count = 1
	pinger.Timeout = time.Second
	pinger.SetPrivileged(true)

	err = pinger.Run()
	if err != nil {
		pingLog.Success = false
	} else {
		pingLog.Success = true
	}
	pingLog.Pingtime = pinger.Statistics().AvgRtt.Seconds() * 1000
	return pingLog, nil
}

func (s *Service) GetContainers() []models.Container {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	containersCopy := make([]models.Container, len(s.containers))
	copy(containersCopy, s.containers)

	return containersCopy
}
