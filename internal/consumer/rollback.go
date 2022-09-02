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

const HandleRollback = "handle_rollback"

type rollbackConsumer struct {
	cfg config.Config

	repo   storage.PersonRepo
	broker broker.Broker
}

func NewRollbackConsumer(cfg config.Config, broker broker.Broker, repo storage.PersonRepo) *rollbackConsumer {
	return &rollbackConsumer{cfg: cfg, broker: broker, repo: repo}
}

func (c *rollbackConsumer) Consume(ctx context.Context, topic string) error {
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
				logger.Errorf("Consumer error: %v", err)
				break
			}

			if err = c.handle(msg.Ctx, person.ID); err != nil {
				logger.Errorf(err.Error())
			}
		}
	}
}

func (c *rollbackConsumer) handle(ctx context.Context, id uint64) error {
	var err error

	span, ctx := opentracing.StartSpanFromContext(ctx, HandleRollback)
	defer span.Finish()

	defer func() {
		if err != nil {
			message_processed_fail.Inc()
		} else {
			message_processed_ok.Inc()
		}
	}()

	err = c.repo.Delete(ctx, id)

	return err
}
