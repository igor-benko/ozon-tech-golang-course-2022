package storage

import (
	"sync"
	"sync/atomic"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type memoryStorage struct {
	sync.RWMutex

	kv        map[uint64]entity.Person
	currentID uint64
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		kv: make(map[uint64]entity.Person),
	}
}

func (ms *memoryStorage) Create(item entity.Person) (uint64, error) {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.kv[item.ID]; ok {
		return 0, entity.ErrPersonAlreadyExists
	}

	item.ID = ms.nextID()
	ms.kv[item.ID] = item
	return item.ID, nil
}

func (ms *memoryStorage) Update(item entity.Person) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.kv[item.ID]; !ok {
		return entity.ErrPersonNotFound
	}

	ms.kv[item.ID] = item
	return nil
}

func (ms *memoryStorage) Delete(personID uint64) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.kv[personID]; !ok {
		return entity.ErrPersonNotFound
	}

	delete(ms.kv, personID)
	return nil
}

func (ms *memoryStorage) List() []entity.Person {
	ms.RLock()
	defer ms.RUnlock()

	items := make([]entity.Person, 0, len(ms.kv))
	for _, v := range ms.kv {
		items = append(items, v)
	}

	return items
}

func (ms *memoryStorage) nextID() uint64 {
	atomic.AddUint64(&ms.currentID, 1)
	return ms.currentID
}
