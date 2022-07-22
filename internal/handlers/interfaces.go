package handlers

import (
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

// Определяем структуру там где ее используем - делаем слабую связность
type Storage interface {
	Create(item entity.Person) (uint64, error)
	Update(personIDitem entity.Person) error
	Delete(personID uint64) error
	List() []entity.Person
}
