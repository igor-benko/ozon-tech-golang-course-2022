# Запуск
В файле config.yaml укажите телеграм API_KEY

# Поддерживаемый список команд

Создание персоны
/person create {LastName} {FirstName}

Обновление персоны
/person update {personID} {LastName} {FirstName}

Удаление персоны
/person delete {personID}

Список персон
/person list


# Архитектура

- entity - содержит бизнес сущности
- service - сервисы
- handlers - обработчики телеграм команд (одна сущность - один хэндлер)
- storage - хранилище (in-memory)