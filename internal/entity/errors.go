package entity

import "errors"

// Common errors
var ErrTimeout = errors.New("Превышел таймаут ожидания")

// Person errors
var ErrPersonNotFound = errors.New("Персона не найдена")
var ErrPersonAlreadyExists = errors.New("Персона уже существует")

// Vehicle errors
var ErrVehicleNotFound = errors.New("Транспорт не найден")
var ErrVehicleAlreadyExists = errors.New("Транспорт уже существует")
