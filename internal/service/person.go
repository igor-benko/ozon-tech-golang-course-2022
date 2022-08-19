package service

import (
	"context"
	"time"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
	repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type personService struct {
	pb.UnimplementedPersonServiceServer

	person  repo.PersonRepo
	vehicle repo.VehicleRepo
	cfg     config.Config
}

func NewPersonService(person repo.PersonRepo, vehicle repo.VehicleRepo, cfg config.Config) personService {
	return personService{
		person:  person,
		vehicle: vehicle,
		cfg:     cfg,
	}
}

func (s *personService) CreatePerson(ctx context.Context, req *pb.CreatePersonRequest) (*pb.CreatePersonResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	id, err := s.person.Create(ctx, entity.Person{
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

	err := s.person.Update(ctx, entity.Person{
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

	err := s.person.Delete(ctx, req.GetId())
	if err != nil {
		return nil, handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *personService) GetPerson(ctx context.Context, req *pb.GetPersonRequest) (*pb.GetPersonResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	personID := req.GetId()
	person, err := s.person.Get(ctx, personID)
	if err != nil {
		return nil, handleError(err)
	}

	vehicles, err := s.vehicle.GetByPersonID(ctx, personID)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.GetPersonResponse{
		Person: &pb.Person{
			Id:        person.ID,
			LastName:  person.LastName,
			FirstName: person.FirstName,
			Vehicles:  mapVehicleToPbVehicle(vehicles),
		},
	}, nil
}

func (s *personService) ListPerson(req *pb.ListPersonRequest, stream pb.PersonService_ListPersonServer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	page, err := s.person.List(ctx, entity.PersonFilter{
		Offset: req.GetOffset(),
		Limit:  req.GetLimit(),
		Order:  req.GetOrder(),
	})
	if err != nil {
		return handleError(err)
	}

	for _, person := range page.Persons {
		vehicles, err := s.vehicle.GetByPersonID(ctx, person.ID)
		if err != nil {
			return handleError(err)
		}

		err = stream.Send(&pb.Person{
			Id:        person.ID,
			LastName:  person.LastName,
			FirstName: person.FirstName,
			Vehicles:  mapVehicleToPbVehicle(vehicles),
		})
		if err != nil {
			return handleError(err)
		}
	}

	return nil
}
