# Запуск
В файле config.yaml укажите телеграм API_KEY и настройки портов для GRPC и Gateway серверов

Затем вызов "make run-server" - сборка и запуск серверов

# Поддерживаемый список команд

Создание персоны
/person create {LastName} {FirstName}

Обновление персоны
/person update {personID} {LastName} {FirstName}

Удаление персоны
/person delete {personID}

Список персон
/person list


Для GRPC
Создание персоны
POST http://localhost:xxxx/v1/persons
{
    "lastName": "A",
    "firstName": "B"
}

Обновление персоны
PUT http://localhost:xxxx/v1/persons
{
    "id": 1,
    "lastName": "A",
    "firstName": "B"
}

Удаление персоны
DELETE http://localhost:xxxx/v1/persons/{id}

Список персон
GET http://localhost:xxxx/v1/persons

Swagger
http://localhost:xxxx/swagger/index.html

# Архитектура

- entity - содержит бизнес сущности
- service - сервисы
- handlers - обработчики телеграм команд (одна сущность - один хэндлер)
- storage - хранилище (in-memory)