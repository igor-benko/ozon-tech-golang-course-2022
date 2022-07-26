package commander

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
)

type commander struct {
	api    BotAPI
	person CommandHandler
}

func NewCommander(api BotAPI, person CommandHandler) commander {
	return commander{
		api:    api,
		person: person,
	}
}

func (c *commander) Run() {
	cfg := config.Get()

	u := tgbotapi.NewUpdate(cfg.Telegram.Offset)
	u.Timeout = cfg.Telegram.Timeout

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

	switch command {
	case "person":
		outputMessage = handleAction(c.person, args...)
	// В случае появления новых сущностей - добавляем их тут
	// case "order":
	// 	outputMessage = handleAction(c.order, args...)
	default:
		outputMessage = "Неизвестная команда"
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, outputMessage)
	c.api.Send(msg)
}

func handleAction(h CommandHandler, args ...string) string {
	action := args[0]

	switch action {
	case "create":
		return h.Create(args...)
	case "update":
		return h.Update(args...)
	case "delete":
		return h.Delete(args...)
	case "list":
		return h.List(args...)
	}

	return "Неподдерживаемая команда"
}
