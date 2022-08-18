package commander

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/adapters"
)

type personHandler struct {
	service adapters.Person
}

func NewPersonHandler(service adapters.Person) *personHandler {
	return &personHandler{
		service: service,
	}
}

func (c *personHandler) Create(ctx context.Context, args ...string) string {
	if len(args) != 3 {
		return "Неправильный формат. Должно быть /person create фамилия имя"
	}

	id, err := c.service.CreatePerson(ctx, args[1], args[2])

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

	if err = c.service.UpdatePerson(ctx, uint64(id), args[2], args[3]); err != nil {
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

	if err = c.service.DeletePerson(ctx, uint64(id)); err != nil {
		return fmt.Sprintf("Ошибка удаления персоны: %s", err)
	}

	return fmt.Sprintf("Удалена персона с ID: %d", id)
}

func (c *personHandler) List(ctx context.Context, args ...string) string {
	var offset uint64
	if len(args) == 2 {
		value, err := strconv.Atoi(args[1])
		if err != nil {
			return "Неправильный формат offset"
		}
		offset = uint64(value)
	}

	var limit uint64
	if len(args) == 3 {
		value, err := strconv.Atoi(args[2])
		if err != nil {
			return "Неправильный формат limit"
		}
		limit = uint64(value)
	}

	var order string
	if len(args) == 4 {
		order = args[3]
	}

	count := 0
	outputMessage := strings.Builder{}
	dataCh, errCh := c.service.ListPerson(ctx, offset, limit, order)

loop:
	for {
		select {
		case item, ok := <-dataCh:
			if !ok {
				break loop
			}
			outputMessage.WriteString(fmt.Sprintf("%d - %s %s\n", item.ID, item.LastName, item.FirstName))
			count++
		case err, ok := <-errCh:
			if !ok {
				break loop
			}
			return err.Error()
		}
	}

	if count == 0 {
		return "Персон нет"
	}

	return outputMessage.String()
}
