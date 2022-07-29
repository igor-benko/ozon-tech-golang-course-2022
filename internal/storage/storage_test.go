package storage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

var storageConfig = config.StorageConfig{
	PoolSize: 10,
}

func TestStorage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	s := NewMemoryStorage(storageConfig)

	// Тест создания
	id, err := s.Create(ctx, entity.Person{
		LastName:  "A",
		FirstName: "B",
	})

	assert.NoError(t, err)
	assert.Equal(t, id, uint64(1))

	_, err = s.Create(ctx, entity.Person{
		ID:        1,
		LastName:  "A",
		FirstName: "B",
	})

	assert.Error(t, err)

	// Тест обновления
	err = s.Update(ctx, entity.Person{
		ID:        1,
		LastName:  "C",
		FirstName: "D",
	})

	assert.NoError(t, err)

	items, err := s.List(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))

	// Тест удаления
	err = s.Delete(ctx, 1)
	assert.NoError(t, err)

	items, err = s.List(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(items))
}

func TestStorageConcurent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	s := NewMemoryStorage(storageConfig)
	wg := &sync.WaitGroup{}

	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := s.Create(ctx, entity.Person{
				LastName:  "C",
				FirstName: "D",
			})

			assert.NoError(t, err)
		}()
	}

	wg.Wait()

	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			err := s.Update(ctx, entity.Person{
				ID:        uint64(id),
				LastName:  "C",
				FirstName: "D",
			})

			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			items, err := s.List(ctx)
			assert.NoError(t, err)
			assert.Equal(t, 999, len(items))
		}(i)
	}

	wg.Wait()

	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			err := s.Delete(ctx, uint64(id))
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()
}
