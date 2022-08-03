package storage

import (
	"context"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type Storage interface {
	Create(ctx context.Context, item entity.Person) (uint64, error)
	Update(ctx context.Context, personIDitem entity.Person) error
	Delete(ctx context.Context, personID uint64) error
	Get(ctx context.Context, id uint64) (*entity.Person, error)
	List(ctx context.Context) ([]entity.Person, error)
}
