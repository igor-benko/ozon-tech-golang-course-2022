package adapters

import (
	"context"
	"io"

	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
)

type PersonDto struct {
	ID        uint64
	LastName  string
	FirstName string
}

type Person interface {
	CreatePerson(ctx context.Context, lastName, firstName string) (uint64, error)
	UpdatePerson(ctx context.Context, id uint64, lastName, firstName string) error
	DeletePerson(ctx context.Context, id uint64) error
	ListPerson(ctx context.Context, offset, limit uint64, order string) (<-chan *PersonDto, <-chan error)
}

type person struct {
	client pb.PersonServiceClient
}

func NewPerson(client pb.PersonServiceClient) *person {
	return &person{
		client: client,
	}
}

func (s *person) CreatePerson(ctx context.Context, lastName, firstName string) (uint64, error) {
	resp, err := s.client.CreatePerson(ctx, &pb.CreatePersonRequest{
		LastName:  lastName,
		FirstName: firstName,
	})

	return resp.GetId(), err
}

func (s *person) UpdatePerson(ctx context.Context, id uint64, lastName, firstName string) error {
	_, err := s.client.UpdatePerson(ctx, &pb.UpdatePersonRequest{
		Id:        id,
		LastName:  lastName,
		FirstName: firstName,
	})

	return err
}

func (s *person) DeletePerson(ctx context.Context, id uint64) error {
	_, err := s.client.DeletePerson(ctx, &pb.DeletePersonRequest{
		Id: id,
	})

	return err
}

func (s *person) ListPerson(ctx context.Context, offset, limit uint64, order string) (<-chan *PersonDto, <-chan error) {
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
