// go:generate mockgen -source ./repository.go -destination ./repository_mock.go -package=storage
package storage

import (
	"context"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type PersonRepo interface {
	Create(ctx context.Context, item entity.Person) (uint64, error)
	Update(ctx context.Context, item entity.Person) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*entity.Person, error)
	List(ctx context.Context, filter entity.PersonFilter) (*entity.PersonPage, error)
}

type VehicleRepo interface {
	Create(ctx context.Context, item entity.Vehicle) (uint64, error)
	Update(ctx context.Context, item entity.Vehicle) error
	Delete(ctx context.Context, id uint64) error
	Exists(ctx context.Context, regNum string) (bool, error)
	Get(ctx context.Context, id uint64) (*entity.Vehicle, error)
	GetByPersonID(ctx context.Context, personID uint64) ([]entity.Vehicle, error)
	List(ctx context.Context, filter entity.VehicleFilter) (*entity.VehiclePage, error)
}
