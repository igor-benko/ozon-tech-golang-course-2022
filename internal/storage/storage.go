package storage

import (
	"context"
	"sync"
	"sync/atomic"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type memoryStorage struct {
	sync.RWMutex

	kv        map[uint64]entity.Person
	currentID uint64
	pool      chan struct{}
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		kv:   make(map[uint64]entity.Person),
		pool: make(chan struct{}, config.Get().Storage.PoolSize),
	}
}

func (ms *memoryStorage) Create(ctx context.Context, item entity.Person) (uint64, error) {
	result, err := ms.wrap(ctx, func() (interface{}, error) {
		if _, ok := ms.kv[item.ID]; ok {
			return 0, entity.ErrPersonAlreadyExists
		}

		item.ID = ms.nextID()
		ms.kv[item.ID] = item
		return item.ID, nil
	})

	if err != nil {
		return 0, err
	}

	return result.(uint64), nil
}

func (ms *memoryStorage) Update(ctx context.Context, item entity.Person) error {
	_, err := ms.wrap(ctx, func() (interface{}, error) {
		if _, ok := ms.kv[item.ID]; !ok {
			return nil, entity.ErrPersonNotFound
		}

		ms.kv[item.ID] = item
		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (ms *memoryStorage) Delete(ctx context.Context, personID uint64) error {
	_, err := ms.wrap(ctx, func() (interface{}, error) {
		if _, ok := ms.kv[personID]; !ok {
			return nil, entity.ErrPersonNotFound
		}

		delete(ms.kv, personID)
		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (ms *memoryStorage) Get(ctx context.Context, personID uint64) (*entity.Person, error) {
	result, err := ms.wrap(ctx, func() (interface{}, error) {
		person, ok := ms.kv[personID]
		if !ok {
			return nil, entity.ErrPersonNotFound
		}

		// Копируем данные персоны чтобы нельзя было изменить ее по адресу
		return &entity.Person{
			ID:        person.ID,
			LastName:  person.LastName,
			FirstName: person.FirstName,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*entity.Person), nil

}

func (ms *memoryStorage) List(ctx context.Context) ([]entity.Person, error) {
	result, err := ms.wrap(ctx, func() (interface{}, error) {
		items := make([]entity.Person, 0, len(ms.kv))
		for _, v := range ms.kv {
			items = append(items, v)
		}

		return items, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]entity.Person), nil

}

func (ms *memoryStorage) nextID() uint64 {
	atomic.AddUint64(&ms.currentID, 1)
	return ms.currentID
}

type WrapFunc func() (interface{}, error)

func (ms *memoryStorage) wrap(ctx context.Context, f WrapFunc) (interface{}, error) {
	resultCh := make(chan interface{}, 1)
	errorCh := make(chan error, 1)

	go func() {
		ms.pool <- struct{}{}
		ms.Lock()

		defer func() {
			ms.Unlock()
			<-ms.pool
		}()

		result, err := f()
		if err != nil {
			errorCh <- err
			return
		}

		resultCh <- result
	}()

	select {
	case <-ctx.Done():
		return nil, entity.ErrTimeout
	case result := <-resultCh:
		return result, nil
	case err := <-errorCh:
		return nil, err
	}
}
