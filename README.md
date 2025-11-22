**О Проекте**
Сервис для автоматического назначения ревьюверов на Pull Request'ы внутри команд разработки.

**Функционал**
Назначение ревьюверов из команды автора pr.

Запрет изменений после merge pr.

Назначение только активных пользователей (isActive = true).

Ограничение до 2 ревьюверов.

Идемпотентность операции merge.

_стек_:
Go, Docker/Docker Compose, in-memory,порт 8080

**Запуск проекта**

Сервис запускается через Docker Compose:

Запуск
`docker-compose up`

Локальный запуск для разработки
`make run`

Сборка проекта
`make build`

Очистка артефактов
`make clean`

Запуск docker-compose
`docker-up`

Конфигурация реализована через переменные окружения в docker-compose.yml:

```yaml
environment:
  - PORT=8080
```

Порт можно изменить при необходимости, обновив значение в docker-compose.yml.

---

_Endpoints_
Все эндпоинты запускаются по стандартному маршруту `localhost:8080`

- Проверка работы сервера
  GET `http://localhost:8080/health - {"status":"ok"}`

- Создание команды POST `http://localhost:8080/team/add`

  ```json
  {
  	"team_name": "backend",
  	"members": [
  		{ "user_id": "1", "username": "Ramz", "is_active": true },
  		{ "user_id": "2", "username": "Ars", "is_active": true },
  		{ "user_id": "3", "username": "Artem", "is_active": false }
  	]
  }
  ```

- Информация о команде
  GET `http://localhost:8080/team/get?team_name=backend`

- Акт/деакт юзера
  POST `http://localhost:8080/users/setIsActive`

  ```json
  {
  	"user_id": "3",
  	"is_active": true
  }
  ```

- Создание pull request
  POST `http://localhost:8080/pullRequest/create`

  ```json
  {
  	"pull_request_id": "pr-1001",
  	"pull_request_name": "Add search",
  	"author_id": "1"
  }
  ```

- Перераспределение ревьювера
  POST `http://localhost:8080/pullRequest/reassign`

  ```json:
  {
  "pull_request_id": "pr-1001",
  "old_user_id": "2"
  }

  ```

- Мёрж pr
  POST `http://localhost:8080/pullRequest/merge`

  ```json
  {
  	"pull_request_id": "pr-1001"
  }
  ```

- история ревью юзера
  GET `http://localhost:8080/users/getReview?user_id=2`

- Лимит 2 ревьюверов
  POST `http://localhost:8080/pullRequest/create`

```json
{
	"pull_request_id": "pr-test-limit",
	"pull_request_name": "Test Limit",
	"author_id": "u1"
}
```

- Запрет изменений после merge
  POST `http://localhost:8080/pullRequest/merge`

```json
{
	"pull_request_id": "pr-test-limit"
}
```

POST `http://localhost:8080/pullRequest/reassign`

```json
{
	"pull_request_id": "pr-test-limit",
	"old_user_id": "u3"
}
```

- Идемпотентность merge
  POST `http://localhost:8080/pullRequest/merge`

```json
{
	"pull_request_id": "pr-test-limit"
}
```

Примерный ответ:

```json
{
	"pr": {
		"pull_request_id": "pr-test-limit",
		"...": "...",
		"mergedAt": "2025-11-22T17:17:55.08124851+03:00"
	}
}
```

- Структурированные ошибки
  GET `http://localhost:8080/team/get?team_name=nonexistent-team`

```json
{ "error": "NotFound" }
```

HTTP Status:404

_Дополнительно_

Проверка работы сервера в терминале.

Лимит 2 ревьювера:

```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-test-limit",
    "pull_request_name": "Test Limit",
    "author_id": "u1"
  }'
```

Запрет изменений после merge:

```bash
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{"pull_request_id": "pr-test-limit"}'
```

```bash
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{"pull_request_id": "pr-test-limit", "old_user_id": "u3"}'
```

Идемпотентность merge:

```bash
#Второй вызов merge:
{"pr":{"pull_request_id":"pr-test-limit",...,"mergedAt":"2025-11-22T17:17:55.08124851+03:00"}}
```

```bash
# третий вызов merge:
{"pr":{"pull_request_id":"pr-test-limit",...,"mergedAt":"2025-11-22T17:17:55.08124851+03:00"}}
```

Структурированные ошибки(частично):

```bash
curl -X GET "http://localhost:8080/team/get?team_name=nonexistent-team" \
  -w "Status: %{http_code}\n"
```

- [Спецификация API](./api/openapi.yaml)
