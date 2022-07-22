package handlers

import (
	"fmt"
	"strconv"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type personHandler struct {
	storage Storage
}

func NewPersonHandler(storage Storage) *personHandler {
	return &personHandler{
		storage: storage,
	}
}

func (c *personHandler) Create(args ...string) string {
	if len(args) != 3 {
		// msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Неправильный формат. Должно быть /person create фамилия имя")
		// c.api.Send(msg)
		return "Неправильный формат. Должно быть /person create фамилия имя"
	}

	id, err := c.storage.Create(entity.Person{
		LastName:  args[1],
		FirstName: args[2],
	})

	if err != nil {
		return fmt.Sprintf("Ошибка создания персоны: %s", err)
	}

	return fmt.Sprintf("Создана персона с ID: %d", id)
}

func (c *personHandler) Update(args ...string) string {
	if len(args) != 4 {
		return "Неправильный формат. Должно быть /person update id фамилия имя"
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		return "Неправильный формат идентификатора"
	}

	err = c.storage.Update(entity.Person{
		ID:        uint64(id),
		LastName:  args[2],
		FirstName: args[3],
	})

	if err != nil {
		return fmt.Sprintf("Ошибка создания персоны: %s", err)
	}

	return fmt.Sprintf("Обновлена персона с ID: %d", id)
}

func (c *personHandler) Delete(args ...string) string {
	if len(args) != 2 {
		return "Неправильный формат. Должно быть /person delete id"
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		return "Неправильный формат идентификатора"
	}

	err = c.storage.Delete(uint64(id))

	if err != nil {
		return fmt.Sprintf("Ошибка удаления персоны: %s", err)
	}

	return fmt.Sprintf("Удалена персона с ID: %d", id)
}

func (c *personHandler) List(args ...string) string {
	outputMessage := ""
	items := c.storage.List()
	if len(items) == 0 {
		outputMessage = "Персон нет"
	} else {
		for _, item := range items {
			outputMessage += fmt.Sprintf("%d - %s %s", item.ID, item.LastName, item.FirstName) + "\n"
		}
	}

	return outputMessage
}
