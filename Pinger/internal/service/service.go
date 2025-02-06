package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pinger/internal/config"
	"pinger/internal/models"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type Service struct {
	client     *http.Client
	containers []models.Container
	mutex      sync.Mutex
}

func NewService() *Service {
	return &Service{
		client: &http.Client{Timeout: config.HTTPTimeout},
	}
}

func (s *Service) FetchContainers() error {
	resp, err := s.client.Get(config.ServerURL + "/containers")
	if err != nil {
		return fmt.Errorf("failed to fetch containers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected server response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var containers []models.Container
	if err := json.Unmarshal(body, &containers); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	s.mutex.Lock()
	s.containers = containers
	s.mutex.Unlock()

	return nil
}

func (s *Service) SendPingLog(pingLog models.PostPingLog) error {
	jsonPingLog, err := json.Marshal(pingLog)
	if err != nil {
		return fmt.Errorf("error serializing JSON: %w", err)
	}

	req, err := http.NewRequest("POST", config.ServerURL+"/pinglog", bytes.NewBuffer(jsonPingLog))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server returned unexpected status: %d", resp.StatusCode)
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

	return pingLog, nil
}

func (s *Service) GetContainers() []models.Container {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	containersCopy := make([]models.Container, len(s.containers))
	copy(containersCopy, s.containers)

	return containersCopy
}
