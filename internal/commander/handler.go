package commander

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/cache"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/client"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

const (
	RequestType_PUBSUB = 0
	RequestType_RETRY  = 1
)

type personHandler struct {
	client client.PersonClient
	cache  cache.CacheClient
}

func NewPersonHandler(client client.PersonClient, cache cache.CacheClient) *personHandler {
	return &personHandler{
		client: client,
		cache:  cache,
	}
}

func (c *personHandler) Create(ctx context.Context, args ...string) string {
	if len(args) != 3 {
		return "Неправильный формат. Должно быть /person create фамилия имя"
	}

	err := c.client.CreatePerson(ctx, args[1], args[2])

	if err != nil {
		return fmt.Sprintf("Ошибка создания персоны: %s", err)
	}

	return "Персона создана"
}

func (c *personHandler) Update(ctx context.Context, args ...string) string {
	if len(args) != 4 {
		return "Неправильный формат. Должно быть /person update id фамилия имя"
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		return "Неправильный формат идентификатора"
	}

	if err = c.client.UpdatePerson(ctx, uint64(id), args[2], args[3]); err != nil {
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

	if err = c.client.DeletePerson(ctx, uint64(id)); err != nil {
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
	dataCh, errCh := c.client.ListPerson(ctx, offset, limit, order)

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

func (c *personHandler) ListAll(ctx context.Context, args ...string) string {
	if len(args) != 2 {
		return "Неправильный формат. Должно быть /person list_all requestType"
	}

	requestType, err := strconv.Atoi(args[1])
	if err != nil {
		return "Неправильный формат типа запроса"
	}

	key := uuid.NewString()
	resultCh := make(chan entity.Person)
	ready := make(chan struct{})

	go func(key string) {
		ready <- struct{}{}

		result := []entity.Person{}

		if requestType == RequestType_PUBSUB {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			dataCh := c.cache.Subscribe(ctx)
			for item := range dataCh {
				if item.Key == key {
					err = json.Unmarshal(item.Payload, &result)
					if err != nil {
						logger.Warnf("Error while getting data from cache: %s", err.Error())
					}
					break
				}
			}
		} else if requestType == RequestType_RETRY {
			attempts := 3 //TODO: ынести в конфиг
			for attempts >= 0 {
				data, err := c.cache.Get(ctx, key)
				if err == nil {
					err = json.Unmarshal(data, &result)
					if err == nil {
						break
					}
				}

				logger.Warnf("Error while getting data from cache: %s", err.Error())

				attempts--
				time.Sleep(1 * time.Second)
			}
		}

		for _, person := range result {
			resultCh <- person
		}
		close(resultCh)
	}(key)

	<-ready

	outputMessage := strings.Builder{}
	if err = c.client.ListAllPersons(ctx, key, uint32(requestType)); err != nil {
		return err.Error()
	}

	for item := range resultCh {
		outputMessage.WriteString(fmt.Sprintf("%d - %s %s\n", item.ID, item.LastName, item.FirstName))
	}

	return outputMessage.String()
}
