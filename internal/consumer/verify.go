package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/broker"
	storage "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

type verifyConsumer struct {
	cfg config.Config

	repo   storage.PersonRepo
	broker broker.Broker
}

func NewVerifyConsumer(cfg config.Config, broker broker.Broker, repo storage.PersonRepo) *verifyConsumer {
	return &verifyConsumer{cfg: cfg, broker: broker, repo: repo}
}

func (c *verifyConsumer) Consume(ctx context.Context, topic string) error {
	messages, err := c.broker.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return errors.New("consumer was interrupted")
		case msg, ok := <-messages:
			if !ok {
				return errors.New("channel was closed")
			}
			person := entity.Person{}
			if err := json.Unmarshal(msg.Body, &person); err != nil {
				break
			}
			if err = c.handle(msg.Ctx, person.ID); err != nil {
				logger.Errorf(err.Error())
			}
		}
	}
}

func (c *verifyConsumer) handle(ctx context.Context, id uint64) error {
	var err error

	span, _ := opentracing.StartSpanFromContext(ctx, "handle_verify")
	defer span.Finish()

	defer func() {
		if err != nil {
			message_processed_fail.Inc()
		} else {
			message_processed_ok.Inc()
		}
	}()

	person, err := c.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if person.LastName == "Ivanov" {
		jsonData, err := json.Marshal(&entity.Person{
			ID: id,
		})
		if err != nil {
			return err
		}

		c.broker.Publish(ctx, &broker.Message{
			Topic:  c.cfg.Kafka.ErrorTopic,
			Body:   jsonData,
			Ctx:    ctx,
			Action: "Rollback",
		})
	}

	return err
}
