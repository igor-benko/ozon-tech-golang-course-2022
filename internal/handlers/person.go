package handlers

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type personHandler struct {
	api    BotAPI
	person PersonService
}

func NewPersonHandler(api BotAPI, person PersonService) *personHandler {
	return &personHandler{
		api:    api,
		person: person,
	}
}

func (c *personHandler) HandleCommand(msg *tgbotapi.Message) {
	args := strings.Split(msg.CommandArguments(), " ")
	if len(args) == 0 {
		c.unsupported(msg)
		return
	}

	switch args[0] {
	case "create":
		c.create(msg)
	case "update":
		c.update(msg)
	case "delete":
		c.delete(msg)
	case "list":
		c.list(msg)
	default:
		c.unsupported(msg)
	}
}

func (c *personHandler) create(inputMessage *tgbotapi.Message) {
	args := strings.Split(inputMessage.CommandArguments(), " ")
	if len(args) != 3 {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Неправильный формат. Должно быть /person create фамилия имя")
		c.api.Send(msg)
		return
	}

	id, err := c.person.Create(entity.Person{
		LastName:  args[1],
		FirstName: args[2],
	})

	if err != nil {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Ошибка создания персоны: %s", err))
		c.api.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Создана персона с ID: %d", id))
	c.api.Send(msg)
}

func (c *personHandler) update(inputMessage *tgbotapi.Message) {
	args := strings.Split(inputMessage.CommandArguments(), " ")
	if len(args) != 4 {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Неправильный формат. Должно быть /person update id фамилия имя")
		c.api.Send(msg)
		return
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Неправильный формат идентификатора")
		c.api.Send(msg)
		return
	}

	err = c.person.Update(entity.Person{
		ID:        uint64(id),
		LastName:  args[2],
		FirstName: args[3],
	})

	if err != nil {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Ошибка создания персоны: %s", err))
		c.api.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Обновлена персона с ID: %d", id))
	c.api.Send(msg)
}

func (c *personHandler) delete(inputMessage *tgbotapi.Message) {
	args := strings.Split(inputMessage.CommandArguments(), " ")
	if len(args) != 2 {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Неправильный формат. Должно быть /person delete id")
		c.api.Send(msg)
		return
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Неправильный формат идентификатора")
		c.api.Send(msg)
		return
	}

	err = c.person.Delete(uint64(id))

	if err != nil {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Ошибка удаления персоны: %s", err))
		c.api.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Удалена персона с ID: %d", id))
	c.api.Send(msg)
}

func (c *personHandler) list(inputMessage *tgbotapi.Message) {
	outputMessage := ""
	items := c.person.List()
	if len(items) == 0 {
		outputMessage = "Персон нет"
	} else {
		for _, item := range items {
			outputMessage += fmt.Sprintf("%d - %s %s", item.ID, item.LastName, item.FirstName) + "\n"
		}
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, outputMessage)
	c.api.Send(msg)
}

func (c *personHandler) unsupported(inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Неизвестная команда")
	c.api.Send(msg)
}
