package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

func TestStorage(t *testing.T) {
	s := NewMemoryStorage()

	// Тест создания
	id, err := s.Create(entity.Person{
		LastName:  "A",
		FirstName: "B",
	})

	assert.NoError(t, err)
	assert.Equal(t, id, uint64(1))

	_, err = s.Create(entity.Person{
		ID:        1,
		LastName:  "A",
		FirstName: "B",
	})

	assert.Error(t, err)

	// Тест обновления
	err = s.Update(entity.Person{
		ID:        1,
		LastName:  "C",
		FirstName: "D",
	})

	assert.NoError(t, err)

	items := s.List()
	assert.Equal(t, 1, len(items))

	// Тест удаления
	err = s.Delete(1)
	items = s.List()

	assert.NoError(t, err)
	assert.Equal(t, 0, len(items))
}

func TestStorageConcurent(t *testing.T) {
	s := NewMemoryStorage()
	wg := &sync.WaitGroup{}

	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := s.Create(entity.Person{
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

			err := s.Update(entity.Person{
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

			items := s.List()
			assert.Equal(t, 999, len(items))
		}(i)
	}

	wg.Wait()

	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			err := s.Delete(uint64(id))
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()
}
