package commander

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
	repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
)

type personHandler struct {
	storage repo.PersonRepo
}

func NewPersonHandler(storage repo.PersonRepo) *personHandler {
	return &personHandler{
		storage: storage,
	}
}

func (c *personHandler) Create(ctx context.Context, args ...string) string {
	if len(args) != 3 {
		// msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Неправильный формат. Должно быть /person create фамилия имя")
		// c.api.Send(msg)
		return "Неправильный формат. Должно быть /person create фамилия имя"
	}

	id, err := c.storage.Create(ctx, entity.Person{
		LastName:  args[1],
		FirstName: args[2],
	})

	if err != nil {
		return fmt.Sprintf("Ошибка создания персоны: %s", err)
	}

	return fmt.Sprintf("Создана персона с ID: %d", id)
}

func (c *personHandler) Update(ctx context.Context, args ...string) string {
	if len(args) != 4 {
		return "Неправильный формат. Должно быть /person update id фамилия имя"
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		return "Неправильный формат идентификатора"
	}

	err = c.storage.Update(ctx, entity.Person{
		ID:        uint64(id),
		LastName:  args[2],
		FirstName: args[3],
	})

	if err != nil {
		return fmt.Sprintf("Ошибка создания персоны: %s", err)
	}

	return fmt.Sprintf("Обновлена персона с ID: %d", id)
}

func (c *personHandler) Delete(ctx context.Context, args ...string) string {
	if len(args) != 2 {
		return "Неправильный формат. Должно быть /person delete id"
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		return "Неправильный формат идентификатора"
	}

	err = c.storage.Delete(ctx, uint64(id))

	if err != nil {
		return fmt.Sprintf("Ошибка удаления персоны: %s", err)
	}

	return fmt.Sprintf("Удалена персона с ID: %d", id)
}

func (c *personHandler) List(ctx context.Context, args ...string) string {
	outputMessage := ""
	page, err := c.storage.List(ctx, entity.PersonFilter{})
	if err != nil {
		return err.Error()
	}

	if len(page.Persons) == 0 {
		outputMessage = "Персон нет"
	} else {
		for _, item := range page.Persons {
			outputMessage += fmt.Sprintf("%d - %s %s", item.ID, item.LastName, item.FirstName) + "\n"
		}
	}

	return outputMessage
}
