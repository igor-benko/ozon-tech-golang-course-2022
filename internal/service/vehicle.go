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

type vehicleService struct {
	pb.UnimplementedVehicleServiceServer

	vehicle repo.VehicleRepo
	cfg     config.Config
}

func NewVehicleService(vehicle repo.VehicleRepo, cfg config.Config) vehicleService {
	return vehicleService{
		vehicle: vehicle,
		cfg:     cfg,
	}
}

func (s *vehicleService) CreateVehicle(ctx context.Context, req *pb.CreateVehicleRequest) (*pb.CreateVehicleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	exists, err := s.vehicle.Exists(ctx, req.GetRegNumber())
	if err != nil {
		return nil, handleError(err)
	}

	if exists {
		return nil, handleError(entity.ErrVehicleAlreadyExists)
	}

	id, err := s.vehicle.Create(ctx, entity.Vehicle{
		Brand:     req.GetBrand(),
		Model:     req.GetModel(),
		RegNumber: req.GetRegNumber(),
		PersonID:  req.GetPersonId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.CreateVehicleResponse{
		Id: id,
	}, nil
}

func (s *vehicleService) UpdateVehicle(ctx context.Context, req *pb.UpdateVehicleRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	exists, err := s.vehicle.Exists(ctx, req.GetRegNumber())
	if err != nil {
		return nil, handleError(err)
	}

	if !exists {
		return nil, handleError(entity.ErrVehicleNotFound)
	}

	err = s.vehicle.Update(ctx, entity.Vehicle{
		ID:        req.GetId(),
		Brand:     req.GetBrand(),
		Model:     req.GetModel(),
		RegNumber: req.GetRegNumber(),
		PersonID:  req.GetPersonId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *vehicleService) DeleteVehicle(ctx context.Context, req *pb.DeleteVehicleRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	err := s.vehicle.Delete(ctx, req.GetId())
	if err != nil {
		return nil, handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *vehicleService) GetVehicle(ctx context.Context, req *pb.GetVehicleRequest) (*pb.GetVehicleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	vehicle, err := s.vehicle.Get(ctx, req.GetId())
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.GetVehicleResponse{
		Vehicle: &pb.Vehicle{
			Id:        vehicle.ID,
			Brand:     vehicle.Brand,
			Model:     vehicle.Model,
			RegNumber: vehicle.RegNumber,
			PersonId:  vehicle.PersonID,
		},
	}, nil
}

func (s *vehicleService) ListVehicle(ctx context.Context, req *pb.ListVehicleRequest) (*pb.ListVehicleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	page, err := s.vehicle.List(ctx, entity.VehicleFilter{
		Offset: req.GetOffset(),
		Limit:  req.GetLimit(),
		Order:  req.GetOrder(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.ListVehicleResponse{
		Vehicles: mapVehicleToPbVehicle(page.Vehicles),
		Total:    page.Total,
	}, nil
}

func mapVehicleToPbVehicle(vehicles []entity.Vehicle) []*pb.Vehicle {
	pbVehicles := make([]*pb.Vehicle, len(vehicles))
	for i, vehicle := range vehicles {
		pbVehicles[i] = &pb.Vehicle{
			Id:        vehicle.ID,
			Brand:     vehicle.Brand,
			Model:     vehicle.Model,
			RegNumber: vehicle.RegNumber,
			PersonId:  vehicle.PersonID,
		}
	}

	return pbVehicles
}
