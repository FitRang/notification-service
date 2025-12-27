package eventbus

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	c *kafka.Consumer
}

func (e *EventBus) NewConsumer(groupID string, topics []string) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": e.cfg.Brokers,

		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     e.cfg.Username,
		"sasl.password":     e.cfg.Password,

		"group.id":          groupID,
		"auto.offset.reset": "earliest",

		"enable.auto.commit": false,
	})

	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	consumer := &Consumer{c: c}
	consumer.startDeliveryLoop()

	return consumer, nil
}

func (p *Consumer) startDeliveryLoop() {
	go func() {
		for ev := range p.c.Events() {
			switch e := ev.(type) {

			case *kafka.Message:
				if e.TopicPartition.Error != nil {
					log.Printf("consumer error: %v", e.TopicPartition.Error)
					continue
				}

				log.Printf(
					"received message topic=%s partition=%d offset=%d",
					*e.TopicPartition.Topic,
					e.TopicPartition.Partition,
					e.TopicPartition.Offset,
				)

				// TODO: call handler here
				// handler(e.Value)

				_, err := p.c.CommitMessage(e)
				if err != nil {
					log.Printf("commit failed: %v", err)
				}

			case kafka.Error:
				log.Printf("kafka error: %v", e)
			}
		}
	}()
}

func (p *Consumer) Close() {
	if err := p.c.Close(); err != nil {
		log.Printf("error closing kafka consumer: %v", err)
	}
}
