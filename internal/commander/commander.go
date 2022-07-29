package commander

import (
	"context"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
)

type commander struct {
	api    BotAPI
	person CommandHandler
	cfg    config.Config
}

func NewCommander(api BotAPI, person CommandHandler, cfg config.Config) commander {
	return commander{
		api:    api,
		person: person,
		cfg:    cfg,
	}
}

func (c *commander) Run() {
	u := tgbotapi.NewUpdate(c.cfg.Telegram.Offset)
	u.Timeout = c.cfg.Telegram.Timeout

	updates := c.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			c.handleCommand(update.Message)
		}
	}
}

func (c *commander) Stop() {
	c.api.StopReceivingUpdates()
}

func (c *commander) handleCommand(inputMessage *tgbotapi.Message) {
	outputMessage := ""

	// Подготовим вход для handler
	command := inputMessage.Command()
	args := strings.Split(inputMessage.CommandArguments(), " ")
	if len(args) == 0 {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Формат команды должен быть /{command} {action} {params}")
		c.api.Send(msg)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cfg.Storage.TimeoutMs)*time.Millisecond)
	defer cancel()

	switch command {
	case "person":
		outputMessage = handleAction(ctx, c.person, args...)
	// В случае появления новых сущностей - добавляем их тут
	// case "order":
	// 	outputMessage = handleAction(c.order, args...)
	default:
		outputMessage = "Неизвестная команда"
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, outputMessage)
	c.api.Send(msg)
}

func handleAction(ctx context.Context, h CommandHandler, args ...string) string {
	action := args[0]

	switch action {
	case "create":
		return h.Create(ctx, args...)
	case "update":
		return h.Update(ctx, args...)
	case "delete":
		return h.Delete(ctx, args...)
	case "list":
		return h.List(ctx, args...)
	}

	return "Неподдерживаемая команда"
}
