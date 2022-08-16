package commander

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
)

type personHandler struct {
	service pb.PersonServiceClient
}

func NewPersonHandler(service pb.PersonServiceClient) *personHandler {
	return &personHandler{
		service: service,
	}
}

func (c *personHandler) Create(ctx context.Context, args ...string) string {
	if len(args) != 3 {
		return "Неправильный формат. Должно быть /person create фамилия имя"
	}

	resp, err := c.service.CreatePerson(ctx, &pb.CreatePersonRequest{
		LastName:  args[1],
		FirstName: args[2],
	})

	if err != nil {
		return fmt.Sprintf("Ошибка создания персоны: %s", err)
	}

	return fmt.Sprintf("Создана персона с ID: %d", resp.GetId())
}

func (c *personHandler) Update(ctx context.Context, args ...string) string {
	if len(args) != 4 {
		return "Неправильный формат. Должно быть /person update id фамилия имя"
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		return "Неправильный формат идентификатора"
	}

	_, err = c.service.UpdatePerson(ctx, &pb.UpdatePersonRequest{
		Id:        uint64(id),
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

	_, err = c.service.DeletePerson(ctx, &pb.DeletePersonRequest{
		Id: uint64(id),
	})

	if err != nil {
		return fmt.Sprintf("Ошибка удаления персоны: %s", err)
	}

	return fmt.Sprintf("Удалена персона с ID: %d", id)
}

func (c *personHandler) List(ctx context.Context, args ...string) string {
	req := &pb.ListPersonRequest{}
	if len(args) == 2 {
		offset, err := strconv.Atoi(args[1])
		if err != nil {
			return "Неправильный формат offset"
		}
		req.Offset = uint64(offset)
	}

	if len(args) == 3 {
		limit, err := strconv.Atoi(args[2])
		if err != nil {
			return "Неправильный формат offset"
		}
		req.Limit = uint64(limit)
	}

	outputMessage := strings.Builder{}
	stream, err := c.service.ListPerson(ctx, req)
	if err != nil {
		return err.Error()
	}

	count := 0
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err.Error()
		}

		count++
		outputMessage.WriteString(fmt.Sprintf("%d - %s %s\n", item.GetId(), item.GetLastName(), item.GetFirstName()))
	}

	if count == 0 {
		return "Персон нет"
	}

	return outputMessage.String()
}
