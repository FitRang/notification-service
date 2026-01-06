package eventbus

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Config struct {
	Brokers  string
}

type EventBus struct {
	cfg Config
}

func NewEventBus(cfg Config) (*EventBus, error) {
	_, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"security.protocol": "PLAINTEXT",
	})

	if err != nil {
		return nil, err
	}

	return &EventBus{cfg: cfg}, nil
}
