package kafka

import (
	"github.com/Shopify/sarama"
)

type Producer interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}
