package commander

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Клиент к телеге.
type BotAPI interface {
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	StopReceivingUpdates()
}

type CommandHandler interface {
	Create(ctx context.Context, args ...string) string
	Update(ctx context.Context, args ...string) string
	Delete(ctx context.Context, args ...string) string
	List(ctx context.Context, args ...string) string
}
