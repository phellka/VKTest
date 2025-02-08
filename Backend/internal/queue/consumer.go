package queue

import (
	"backend/internal/models"
	"backend/internal/service"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type QueueConsumer struct {
	service               *service.Service
	amqpChannelContainers *amqp.Channel
	amqpChannelPingLogs   *amqp.Channel
}

func NewQueueConsumer(s *service.Service, amqpChannelContainers *amqp.Channel, amqpChannelPingLogs *amqp.Channel) *QueueConsumer {
	return &QueueConsumer{
		service:               s,
		amqpChannelContainers: amqpChannelContainers,
		amqpChannelPingLogs:   amqpChannelPingLogs,
	}
}

func (qc *QueueConsumer) ListenContainersRequest() {
	requestQueue, err := qc.amqpChannelContainers.QueueDeclare(
		os.Getenv("ContainersRequestQueueName"),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %w", err)
	}
	msgs, err := qc.amqpChannelContainers.Consume(
		requestQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to queue %s: %v", requestQueue.Name, err)
	}

	for msg := range msgs {
		log.Println("Received a request for a list of containers")

		containers, err := qc.service.GetContainers()
		if err != nil {
			log.Printf("Error retrieving containers: %v", err)
			continue
		}

		responseBody, err := json.Marshal(containers)
		if err != nil {
			log.Printf("Error serializing containers: %v", err)
			continue
		}
		log.Println(msg.ReplyTo)
		err = qc.amqpChannelContainers.Publish(
			"",
			msg.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/json",
				Body:          responseBody,
				CorrelationId: msg.CorrelationId,
			},
		)
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
}

func (qc *QueueConsumer) ListenPingLogs() {
	replyQueue, err := qc.amqpChannelPingLogs.QueueDeclare(
		os.Getenv("PingLogsQueueName"),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %w", err)
	}
	msgs, err := qc.amqpChannelPingLogs.Consume(
		replyQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to queue %s: %v", replyQueue.Name, err)
	}

	for msg := range msgs {
		var pingLog models.PingLog
		if err := json.Unmarshal(msg.Body, &pingLog); err != nil {
			log.Printf("Error deserializing pingLog: %v", err)
			continue
		}
		err := qc.service.CreatePingLog(pingLog)
		if err != nil {
			log.Printf("Error saving pingLog: %v", err)
		}
	}
}
