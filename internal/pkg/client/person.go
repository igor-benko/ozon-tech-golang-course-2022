package client

import (
	"context"
	"encoding/json"
	"io"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/broker"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
)

type PersonDto struct {
	ID        uint64
	LastName  string
	FirstName string
}

type PersonClient interface {
	CreatePerson(ctx context.Context, lastName, firstName string) error
	UpdatePerson(ctx context.Context, id uint64, lastName, firstName string) error
	DeletePerson(ctx context.Context, id uint64) error
	ListPerson(ctx context.Context, offset, limit uint64, order string) (<-chan *PersonDto, <-chan error)
}

type personClient struct {
	cfg config.Config

	client pb.PersonServiceClient
	broker broker.Broker
}

func NewPersonClient(cfg config.Config, client pb.PersonServiceClient, broker broker.Broker) *personClient {
	return &personClient{
		cfg:    cfg,
		client: client,
		broker: broker,
	}
}

func (s *personClient) CreatePerson(ctx context.Context, lastName, firstName string) error {
	jsonData, err := json.Marshal(&PersonDto{
		LastName:  lastName,
		FirstName: firstName,
	})
	if err != nil {
		return err
	}

	return s.broker.Publish(ctx, &broker.Message{
		Topic:  s.cfg.Kafka.IncomeTopic,
		Action: "Create",
		Body:   jsonData,
	})
}

func (s *personClient) UpdatePerson(ctx context.Context, id uint64, lastName, firstName string) error {
	jsonData, err := json.Marshal(&PersonDto{
		ID:        id,
		LastName:  lastName,
		FirstName: firstName,
	})
	if err != nil {
		return err
	}

	return s.broker.Publish(ctx, &broker.Message{
		Topic:  s.cfg.Kafka.IncomeTopic,
		Action: "Update",
		Body:   jsonData,
	})
}

func (s *personClient) DeletePerson(ctx context.Context, id uint64) error {
	jsonData, err := json.Marshal(&PersonDto{
		ID: id,
	})
	if err != nil {
		return err
	}

	return s.broker.Publish(ctx, &broker.Message{
		Topic:  s.cfg.Kafka.IncomeTopic,
		Action: "Delete",
		Body:   jsonData,
	})
}

func (s *personClient) ListPerson(ctx context.Context, offset, limit uint64, order string) (<-chan *PersonDto, <-chan error) {
	ch := make(chan *PersonDto, 1)
	errCh := make(chan error, 1)

	go func() {
		stream, err := s.client.ListPerson(ctx, &pb.ListPersonRequest{
			Offset: offset,
			Limit:  limit,
			Order:  order,
		})

		if err != nil {
			errCh <- err
			close(ch)
			close(errCh)
			return
		}

		for {
			item, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				errCh <- err
				break
			}

			ch <- &PersonDto{
				ID:        item.GetId(),
				LastName:  item.GetLastName(),
				FirstName: item.GetFirstName(),
			}
		}

		close(ch)
		close(errCh)
	}()

	return ch, errCh
}
