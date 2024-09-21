package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"live-chat-server/config"
	"live-chat-server/internal/mq/types"
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

func (k *kafkaClient) Subscribe(topic string) error {
	if err := k.consumer.Subscribe(topic, nil); err != nil {
		return err
	}
	return nil
}

func (k *kafkaClient) Poll(timeoutMs int) types.Event {
	ev := k.consumer.Poll(timeoutMs)
	switch event := ev.(type) {
	case *kafka.Message:
		return &types.Message{Value: event.Value}
	case *kafka.Error:
		return &types.Error{Error: event}
	default:
		return nil
	}
}
