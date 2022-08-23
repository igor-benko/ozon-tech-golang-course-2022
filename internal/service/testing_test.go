package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	storage "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

type personServiceFixture struct {
	Ctx context.Context

	ctrl        *gomock.Controller
	personRepo  *storage.MockPersonRepo
	vehicleRepo *storage.MockVehicleRepo

	service personService
}

func setUpPersonServiceFixture(t *testing.T) personServiceFixture {
	t.Parallel()

	cfg, err := config.Init()
	if err != nil {
		logger.FatalKV(err.Error())
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	personRepo := storage.NewMockPersonRepo(ctrl)
	vehicleRepo := storage.NewMockVehicleRepo(ctrl)

	service := NewPersonService(personRepo, vehicleRepo, *cfg)

	return personServiceFixture{
		Ctx: ctx,

		ctrl:        ctrl,
		personRepo:  personRepo,
		vehicleRepo: vehicleRepo,

		service: service,
	}
}
