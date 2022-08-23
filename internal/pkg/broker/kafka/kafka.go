package kafka

import (
	"context"
	"time"

	"github.com/Shopify/sarama"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/broker"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

type kafkaBroker struct {
	producer Producer
	consumer Consumer
}

func NewKafkaBroker(producer Producer, consumer Consumer) *kafkaBroker {
	return &kafkaBroker{
		producer: producer,
		consumer: consumer,
	}
}

func (b *kafkaBroker) Publish(ctx context.Context, msg *broker.Message) error {
	var err error

	defer func() {
		if err != nil {
			send_to_kafka_fail.Inc()
		} else {
			send_to_kafka_ok.Inc()
		}
	}()

	_, _, err = b.producer.SendMessage(&sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: sarama.ByteEncoder(msg.Body),
		Headers: []sarama.RecordHeader{
			sarama.RecordHeader{
				Key:   []byte("Action"),
				Value: []byte(msg.Action),
			},
		},
	})

	return err
}

func (b *kafkaBroker) Subscribe(ctx context.Context, topic string) (<-chan *broker.Message, error) {
	messagesCh := make(chan *broker.Message, 1)
	handler := &consumerHandler{messagesCh}

	go func() {
		for {
			if err := b.consumer.Consume(ctx, []string{topic}, handler); err != nil {
				logger.Errorf("on consume: %v", err)
				time.Sleep(time.Second * 10)
			}
		}
	}()

	return messagesCh, nil
}
