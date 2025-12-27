package eventbus

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Config struct {
	Brokers  string
	Username string
	Password string
}

type EventBus struct {
	cfg   Config
	admin *kafka.AdminClient
}

func NewEventBus(cfg Config) (*EventBus, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
	}

	if cfg.Username != "" && cfg.Password != "" {
		conf.SetKey("security.protocol", "SASL_SSL")
		conf.SetKey("sasl.mechanisms", "PLAIN")
		conf.SetKey("sasl.username", cfg.Username)
		conf.SetKey("sasl.password", cfg.Password)
	}

	admin, err := kafka.NewAdminClient(conf)
	if err != nil {
		return nil, err
	}

	return &EventBus{
		cfg:   cfg,
		admin: admin,
	}, nil
}

func (e *EventBus) Close() {
	if e.admin != nil {
		e.admin.Close()
	}
}
