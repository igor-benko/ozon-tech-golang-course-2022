package storage

import (
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type memoryStorage struct {
	kv        map[uint64]entity.Person
	currentID uint64
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		kv: make(map[uint64]entity.Person),
	}
}

func (ms *memoryStorage) Create(item entity.Person) (uint64, error) {
	if _, ok := ms.kv[item.ID]; ok {
		return 0, entity.ErrPersonAlreadyExists
	}

	ms.currentID++
	item.ID = ms.currentID
	ms.kv[ms.currentID] = item
	return ms.currentID, nil
}

func (ms *memoryStorage) Update(item entity.Person) error {
	if _, ok := ms.kv[item.ID]; !ok {
		return entity.ErrPersonNotFound
	}

	ms.kv[item.ID] = item
	return nil
}

func (ms *memoryStorage) Delete(personID uint64) error {
	if _, ok := ms.kv[personID]; !ok {
		return entity.ErrPersonNotFound
	}

	delete(ms.kv, personID)
	return nil
}

func (ms *memoryStorage) List() []entity.Person {
	items := make([]entity.Person, 0, len(ms.kv))
	for _, v := range ms.kv {
		items = append(items, v)
	}

	return items
}
