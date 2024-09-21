package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"live-chat-server/config"
)

type kafkaClient struct {
	cfg      config.Kafka
	consumer *kafka.Consumer
}

func NewKafkaClient(cfg config.Kafka) (Client, error) {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.URL,
		"group.id":          cfg.GroupID,
		"auto.offset.reset": "latest",
	})
	if err != nil {
		return nil, err
	}

	return &kafkaClient{
		cfg:      cfg,
		consumer: consumer,
	}, nil
}
