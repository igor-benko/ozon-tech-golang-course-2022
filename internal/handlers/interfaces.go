package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

// Клиент к телеге. (Лучше маленькое копирование чем новая зависимость)
type BotAPI interface {
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	StopReceivingUpdates()
}

// Сервис по работе с персонами
type PersonService interface {
	Create(item entity.Person) (uint64, error)
	Update(item entity.Person) error
	Delete(personID uint64) error
	List() []entity.Person
}
