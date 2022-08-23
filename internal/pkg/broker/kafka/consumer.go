package kafka

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/broker"
)

type Consumer interface {
	Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error
}

type consumerHandler struct {
	messagesCh chan *broker.Message
}

func (c *consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			log.Print("Done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Print("Data channel closed")
				return nil
			}

			got_from_kafka.Inc()

			action := ""
			for _, h := range msg.Headers {
				if string(h.Key) == "Action" {
					action = string(h.Value)
					break
				}
			}

			c.messagesCh <- &broker.Message{
				Topic:  msg.Topic,
				Action: action,
				Body:   msg.Value,
			}
			session.MarkMessage(msg, "")
		}
	}
}
