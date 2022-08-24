package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/broker"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
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
			logger.Infof("Done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				logger.Infof("Data channel closed")
				return nil
			}

			got_from_kafka.Inc()

			action := ""
			headers := make(map[string]string)
			for _, h := range msg.Headers {
				if string(h.Key) == "Action" {
					action = string(h.Value)
				} else {
					headers[string(h.Key)] = string(h.Value)
				}
			}

			spanContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(headers))
			if err != nil {
				logger.Errorf(err.Error())
			}

			span := opentracing.StartSpan(
				action,
				ext.RPCServerOption(spanContext))

			c.messagesCh <- &broker.Message{
				Topic:  msg.Topic,
				Action: action,
				Body:   msg.Value,
				Ctx:    opentracing.ContextWithSpan(context.Background(), span),
			}
			session.MarkMessage(msg, "")

			defer span.Finish()
		}
	}
}
