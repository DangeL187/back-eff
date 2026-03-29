# Как запустить

Поднять PostgreSQL и приложение (из корня проекта):
```bash
cd docker
docker compose up -d
```

[Для тестов] Поднять тестовую БД на localhost (из корня проекта):
```bash
cd docker/test
docker compose up -d
```

Миграции (golang-migrate) для инициализации БД (в т.ч. тестовой):
```bash
migrate -path ./migrations -database "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable" up
```

Тесты (из корня проекта):
```bash
cd test && go test
```

# Ключевые особенности
- CRUDL эндпоинты
- Update-эндпоинт не имеет side эффекта в виде `INSERT`, так как имеет проверку на существование записи
- Для поиска по датам используются `GIST-индексы`
- `/costs` эндпоинт поддерживает фильтрацию по `user_id` и/или по `service_name`
- Формат дат для `start_date` и `end_date` - `MM-YYYY` согласно ТЗ
- Формат дат для `from` и `to` при запросе к `/costs` - `DD-MM-YYYY`, это даст фильтровать с точностью до дня
- Документация в папке `docs`

# Архитектура
- `DDD` + `Clean Architecture`
  - domain слой
  - handler слой (взаимодействует с usecase)
  - usecase слой (зависит от интерфейсов репозиториев, а не напрямую от БД)
  - infra слой (реализует взаимодействие с БД и внешним миром)
- `Feature-Sliced Design` - моё предпочтение для личных проектов

# Код
- Взаимодействие с БД через `gorm`
- Фреймворк для сервера - `echo`
- error-level логирование в JSON для внутренних ошибок
- debug-level логирование в error-trace string для ошибок по вине пользователя
- Логирование с помощью библиотеки `zap`
- Форматирование ошибок в string и JSON, error-tracing с помощью моей библиотеки [erax](https://github.com/DangeL187/erax). В качестве альтернативы можно было использовать `fmt`
