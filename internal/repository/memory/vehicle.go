package storage

import (
	"context"
	"sync"
	"sync/atomic"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type vehicleStorage struct {
	sync.RWMutex

	kv        map[uint64]entity.Vehicle
	currentID uint64
	pool      chan struct{}
}

func NewVehicleRepo(cfg config.StorageConfig) *vehicleStorage {
	return &vehicleStorage{
		kv:   make(map[uint64]entity.Vehicle),
		pool: make(chan struct{}, cfg.PoolSize),
	}
}

func (ms *vehicleStorage) Create(ctx context.Context, item entity.Vehicle) (uint64, error) {
	result, err := ms.wrap(ctx, func() (interface{}, error) {
		if _, ok := ms.kv[item.ID]; ok {
			return 0, entity.ErrVehicleAlreadyExists
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

func (ms *vehicleStorage) Update(ctx context.Context, item entity.Vehicle) error {
	_, err := ms.wrap(ctx, func() (interface{}, error) {
		if _, ok := ms.kv[item.ID]; !ok {
			return nil, entity.ErrVehicleNotFound
		}

		ms.kv[item.ID] = item
		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (ms *vehicleStorage) Delete(ctx context.Context, id uint64) error {
	_, err := ms.wrap(ctx, func() (interface{}, error) {
		if _, ok := ms.kv[id]; !ok {
			return nil, entity.ErrVehicleNotFound
		}

		delete(ms.kv, id)
		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (ms *vehicleStorage) Get(ctx context.Context, id uint64) (*entity.Vehicle, error) {
	result, err := ms.wrap(ctx, func() (interface{}, error) {
		vehicle, ok := ms.kv[id]
		if !ok {
			return nil, entity.ErrPersonNotFound
		}

		// Копируем данные персоны чтобы нельзя было изменить ее по адресу
		return &entity.Vehicle{
			ID:        vehicle.ID,
			Brand:     vehicle.Brand,
			Model:     vehicle.Model,
			RegNumber: vehicle.RegNumber,
			PersonID:  vehicle.PersonID,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*entity.Vehicle), nil

}

func (ms *vehicleStorage) Exists(ctx context.Context, regNum string) (bool, error) {
	result, err := ms.wrap(ctx, func() (interface{}, error) {
		exists := false
		for _, v := range ms.kv {
			if v.RegNumber == regNum {
				exists = true
				break
			}
		}

		return exists, nil
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (ms *vehicleStorage) GetByPersonID(ctx context.Context, personID uint64) ([]entity.Vehicle, error) {
	result, err := ms.wrap(ctx, func() (interface{}, error) {
		vehicles := []entity.Vehicle{}
		for _, v := range ms.kv {
			if v.PersonID == personID {
				vehicles = append(vehicles, v)
			}
		}

		return vehicles, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]entity.Vehicle), nil
}

func (ms *vehicleStorage) List(ctx context.Context, filter entity.VehicleFilter) (*entity.VehiclePage, error) {
	result, err := ms.wrap(ctx, func() (interface{}, error) {
		items := make([]entity.Vehicle, 0, len(ms.kv))
		for _, v := range ms.kv {
			items = append(items, v)
		}

		return &entity.VehiclePage{
			Vehicles: items,
			Total:    uint64(len(ms.kv)),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*entity.VehiclePage), nil

}

func (ms *vehicleStorage) nextID() uint64 {
	atomic.AddUint64(&ms.currentID, 1)
	return ms.currentID
}

func (ms *vehicleStorage) wrap(ctx context.Context, f WrapFunc) (interface{}, error) {
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
