package eventbus

import (
	"context"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Handler func(key, value []byte) error

type Consumer struct {
	c      *kafka.Consumer
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (e *EventBus) NewConsumer(
	groupID string,
	topics []string,
	handler Handler,
) (*Consumer, error) {

	cfg := &kafka.ConfigMap{
		"bootstrap.servers": e.cfg.Brokers,

		"security.protocol": "PLAINTEXT",
		"sasl.mechanisms":   "PLAIN",

		"group.id":          groupID,
		"auto.offset.reset": "earliest",

		"enable.auto.commit": false,

		"go.application.rebalance.enable": true,
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumer := &Consumer{
		c:      c,
		ctx:    ctx,
		cancel: cancel,
	}

	err = c.SubscribeTopics(topics, func(c *kafka.Consumer, ev kafka.Event) error {
		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			log.Printf("Partitions assigned: %v", e.Partitions)
			return c.Assign(e.Partitions)

		case kafka.RevokedPartitions:
			log.Printf("Partitions revoked: %v", e.Partitions)
			return c.Unassign()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	consumer.start(handler)
	return consumer, nil
}

func (p *Consumer) start(handler Handler) {
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()

		for {
			select {
			case <-p.ctx.Done():
				log.Println("consumer shutting down")
				return

			default:
				ev := p.c.Poll(100)
				if ev == nil {
					continue
				}

				switch e := ev.(type) {

				case *kafka.Message:
					if err := handler(e.Key, e.Value); err != nil {
						log.Printf("handler failed (offset %d): %v", e.TopicPartition.Offset, err)
						continue 
					}

					if _, err := p.c.CommitMessage(e); err != nil {
						log.Printf("commit failed: %v", err)
					}

				case kafka.Error:
					log.Printf("kafka error: %v", e)
				}
			}
		}
	}()
}

func (p *Consumer) Close() {
	p.cancel()
	p.wg.Wait()
	p.c.Close()
}
