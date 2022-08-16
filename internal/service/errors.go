package service

import (
	"errors"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(err error) error {
	if errors.Is(err, entity.ErrPersonAlreadyExists) || errors.Is(err, entity.ErrVehicleAlreadyExists) {
		return status.Error(codes.AlreadyExists, err.Error())
	}

	if errors.Is(err, entity.ErrPersonNotFound) || errors.Is(err, entity.ErrVehicleNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, entity.ErrTimeout) {
		return status.Error(codes.DeadlineExceeded, err.Error())
	}

	return err
}
