package service

import "gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"

type personService struct {
	storage Storage
}

func NewPersonService(storage Storage) *personService {
	return &personService{
		storage: storage,
	}
}

func (s *personService) Create(item entity.Person) (uint64, error) {
	return s.storage.Create(item)
}

func (s *personService) Update(item entity.Person) error {
	return s.storage.Update(item)
}

func (s *personService) Delete(personID uint64) error {
	return s.storage.Delete(personID)
}

func (s *personService) List() []entity.Person {
	return s.storage.List()
}
