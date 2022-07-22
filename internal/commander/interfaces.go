package commander

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Клиент к телеге.
type BotAPI interface {
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	StopReceivingUpdates()
}

type CommandHandler interface {
	Create(args ...string) string
	Update(args ...string) string
	Delete(args ...string) string
	List(args ...string) string
}
