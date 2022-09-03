// go:generate mockgen -source ./../../pkg/api/person_grpc.pb.go -destination ./person_mock.go -package=service
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/cache"
	repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type personService struct {
	pb.UnimplementedPersonServiceServer

	person  repo.PersonRepo
	vehicle repo.VehicleRepo
	cache   cache.CacheClient
	cfg     config.Config
}

func NewPersonService(person repo.PersonRepo, vehicle repo.VehicleRepo, cache cache.CacheClient, cfg config.Config) personService {
	return personService{
		person:  person,
		vehicle: vehicle,
		cache:   cache,
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

	// Get data from cache
	cacheKey := fmt.Sprintf("GetPerson_%d", personID)
	cachedValue, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		response := &pb.GetPersonResponse{}
		err = json.Unmarshal(cachedValue, response)
		if err == nil {
			return response, nil
		}

		logger.Warnf("Failed to unmarshall cachedValue: %s", err.Error())
	}

	// get data from database
	person, err := s.person.Get(ctx, personID)
	if err != nil {
		return nil, handleError(err)
	}

	vehicles, err := s.vehicle.GetByPersonID(ctx, personID)
	if err != nil {
		return nil, handleError(err)
	}

	response := &pb.GetPersonResponse{
		Person: &pb.Person{
			Id:        person.ID,
			LastName:  person.LastName,
			FirstName: person.FirstName,
			Vehicles:  mapVehicleToPbVehicle(vehicles),
		},
	}

	// Save data to cache
	bytes, _ := json.Marshal(response)
	if err = s.cache.Set(ctx, cacheKey, bytes); err != nil {
		logger.Warnf("Error while save data to cache: %s", err.Error())
	}

	return response, nil
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

func (s *personService) ListAllPersons(ctx context.Context, req *pb.ListAllPersonsRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	page, err := s.person.List(ctx, entity.PersonFilter{})
	if err != nil {
		return nil, handleError(err)
	}

	payload, _ := json.Marshal(page.Persons)
	value := &cache.Publication{
		Key:     req.Key,
		Payload: payload,
	}

	if req.RequestType == pb.RequestType_PUBSUB {
		err = s.cache.Publish(ctx, value)
		if err != nil {
			return nil, err
		}
	} else if req.RequestType == pb.RequestType_RETRY {
		err = s.cache.Set(ctx, req.Key, payload)
		if err != nil {
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}
