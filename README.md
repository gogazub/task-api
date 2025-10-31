# Task Service

Сервис для удаленного выполнения пользовательского кода (микро-Colab).  
Клиент отправляет фрагмент кода → сервис ставит задачу в очередь → worker поднимает изолированный Docker-контейнер и исполняет код → статус и результат доступны по API.


---

## Ключевые возможности

- **JWT-аутентификация**: регистрация/логин, защищённые эндпоинты.
- **Swagger/OpenAPI UI**: интерактивное тестирование API в браузере.
- **Асинхронная обработка**: публикация задач в **RabbitMQ**, выполнение в **runner-service**.
- **Изоляция**: исполнение в Docker-контейнерах.
- **Хранение статусов/результатов**: **PostgreSQL**.
- **Health-check**: `GET /` для проверки живости.
- **Docker Compose**: быстрый локальный запуск.

---

## Технологии и паттерны

- **Go** (`net/http`, `http.ServeMux`) + хендлеры:
  - `POST /register`, `POST /login`
  - `POST /task`, `GET /status/{id}`, `GET /result/{id}`
  - `GET /swagger/`, `GET /`
- **JWT Auth** - stateless, с истечением срока действия.
- **OpenAPI/Swagger UI** — автодокументация, тестирование в браузере.
- **Message-Driven Architecture** — RabbitMQ, **async job processing**.
- **Worker pattern + Docker sandbox** — изоляция недоверенного кода.
- **PostgreSQL** — персистентность статусов/результатов.
- **12-factor/containers** — конфигурация через env, `docker-compose`.
- **(Roadmap)** **Redis** для low-latency статусов; **линтеры** (`golangci-lint`) и рефакторинг.

---

## Быстрый старт (Docker Compose)

1) Скопируйте и при необходимости поправьте переменные: .env.example

2) Запустите сервисы через Docker Compose (сборка и старт):
   docker compose up --build

3) Проверьте, что всё поднялось:
   - Health-check: http://localhost:port/  →  "Task Service API is running"
   - Swagger UI:   http://localhost:port/swagger/

---

## API



## Roadmap / TODO

- **Redis** как кэш статусов/результатов (мгновенный `GET /status`).
- **Линтеры/рефакторинг**: `golangci-lint`, интерфейсы, DI, чистая архитектура.
- **Ретраи/идемпотентность** в очереди, дедупликация задач.
- **Rate limiting / CORS**.
- **Метрики и логирование**: Prometheus/OpenTelemetry, структурные логи, трейсинг.
---

