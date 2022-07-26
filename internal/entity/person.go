package entity

import "errors"

type Person struct {
	ID        uint64
	LastName  string
	FirstName string
}

var ErrPersonNotFound = errors.New("Персона не найдена")
var ErrPersonAlreadyExists = errors.New("Персона уже существует")
