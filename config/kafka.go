package config

import "os"

type KafkaConfig struct {
	Brokers  string
	Username string
	Password string
}

func LoadKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers:  os.Getenv("KAFKA_URI"),
		Username: os.Getenv("KAFKA_USERNAME"),
		Password: os.Getenv("KAFKA_PASSWORD"),
	}
}
