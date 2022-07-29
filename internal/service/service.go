package service

import (
	"context"
	"errors"
	"time"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/storage"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type personService struct {
	pb.UnimplementedPersonServiceServer

	storage storage.Storage
	cfg     config.Config
}

func NewPersonService(storage storage.Storage, cfg config.Config) personService {
	return personService{
		storage: storage,
		cfg:     cfg,
	}
}

func (s *personService) CreatePerson(ctx context.Context, req *pb.CreatePersonRequest) (*pb.CreatePersonResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	id, err := s.storage.Create(ctx, entity.Person{
		LastName:  req.GetLastName(),
		FirstName: req.GetFirstName(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.CreatePersonResponse{
		Id: id,
	}, nil
}
func (s *personService) UpdatePerson(ctx context.Context, req *pb.UpdatePersonRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	err := s.storage.Update(ctx, entity.Person{
		ID:        req.GetId(),
		LastName:  req.GetLastName(),
		FirstName: req.GetFirstName(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &emptypb.Empty{}, nil
}
func (s *personService) DeletePerson(ctx context.Context, req *pb.DeletePersonRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	err := s.storage.Delete(ctx, req.GetId())
	if err != nil {
		return nil, handleError(err)
	}

	return &emptypb.Empty{}, nil
}
func (s *personService) GetPerson(ctx context.Context, req *pb.GetPersonRequest) (*pb.GetPersonResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	person, err := s.storage.Get(ctx, req.GetId())
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.GetPersonResponse{
		Person: &pb.Person{
			Id:        person.ID,
			LastName:  person.LastName,
			FirstName: person.FirstName,
		},
	}, nil
}
func (s *personService) ListPerson(ctx context.Context, req *emptypb.Empty) (*pb.ListPersonResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	persons, err := s.storage.List(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	pbPersons := make([]*pb.Person, len(persons))
	for i, person := range persons {
		pbPersons[i] = &pb.Person{
			Id:        person.ID,
			LastName:  person.LastName,
			FirstName: person.FirstName,
		}
	}

	return &pb.ListPersonResponse{
		Persons: pbPersons,
	}, nil
}

func handleError(err error) error {
	if errors.Is(err, entity.ErrPersonAlreadyExists) {
		return status.Error(codes.AlreadyExists, err.Error())
	}

	if errors.Is(err, entity.ErrPersonNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, entity.ErrTimeout) {
		return status.Error(codes.DeadlineExceeded, err.Error())
	}

	return err
}
