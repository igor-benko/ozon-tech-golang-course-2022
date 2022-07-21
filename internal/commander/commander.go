package commander

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			c.HandleCommand(update.Message)
		}
	}
}

func (c *commander) Stop() {
	c.api.StopReceivingUpdates()
}

func (c *commander) HandleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "person":
		c.person.HandleCommand(msg)
	// В случае появления новых сущностей - добавляем их тут
	// case "order":
	// 	c.order.HandleCommand(msg)
	default:
		c.unsupported(msg)
	}
}

func (c *commander) unsupported(inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Неизвестная команда")
	c.api.Send(msg)
}
