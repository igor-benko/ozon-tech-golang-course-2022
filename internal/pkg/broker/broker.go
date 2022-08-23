package broker

import "context"

type Message struct {
	Topic  string
	Action string
	Body   []byte
}

type Broker interface {
	Publish(ctx context.Context, msg *Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)
}
