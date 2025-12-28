# Go HTTP Server

REST API для управления задачами. Написан на стандартной библиотеке.

## Запуск

```bash
# Локально
make run

# Docker
docker-compose up --build
```

## Эндпоинты

- `POST /todos` — создать задачу
- `GET /todos` — список задач
- `GET /todos/{id}` — задача по ID
- `PUT /todos/{id}` — обновить задачу
- `DELETE /todos/{id}` — удалить задачу

---
***Фичи в отдельной ветке***