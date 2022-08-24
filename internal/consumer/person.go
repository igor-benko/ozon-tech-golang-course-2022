package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/broker"
	storage "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

type personConsumer struct {
	cfg config.Config

	repo   storage.PersonRepo
	broker broker.Broker
}

func NewPersonConsumer(cfg config.Config, broker broker.Broker, repo storage.PersonRepo) *personConsumer {
	return &personConsumer{cfg: cfg, broker: broker, repo: repo}
}

func (c *personConsumer) Consume(ctx context.Context, topic string) error {
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

			if err = c.handle(msg.Ctx, msg.Action, person); err != nil {
				logger.Errorf(err.Error())
			}
		}
	}
}

func (c *personConsumer) handle(ctx context.Context, action string, person entity.Person) error {
	var err error

	span, _ := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("handle_%s", action))
	defer span.Finish()

	defer func() {
		if err != nil {
			message_processed_fail.Inc()
		} else {
			message_processed_ok.Inc()
		}
	}()

	switch action {
	case "Create":
		var id uint64
		id, err = c.repo.Create(context.Background(), person)
		if err != nil {
			return err
		}
		err = c.sendToVerify(ctx, id)
	case "Update":
		err = c.repo.Update(context.Background(), person)
	case "Delete":
		err = c.repo.Delete(context.Background(), person.ID)
	}

	return err
}

func (c *personConsumer) sendToVerify(ctx context.Context, id uint64) error {
	jsonData, err := json.Marshal(&entity.Person{
		ID: id,
	})
	if err != nil {
		return err
	}

	err = c.broker.Publish(ctx, &broker.Message{
		Topic:  c.cfg.Kafka.VerifyTopic,
		Body:   jsonData,
		Action: "Verify",
		Ctx:    ctx,
	})

	if err != nil {
		return err
	}

	return nil
}
