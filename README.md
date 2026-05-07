# HTTP Gateway

HTTP Gateway - публичная точка входа для Web/iOS/Android клиентов. Он принимает HTTP-запросы, валидирует protected routes через Auth Service и ходит во внутренние сервисы по gRPC. Kafka напрямую gateway не использует.

## Место в архитектуре

```text
Client
  |
  v
HTTP Gateway
  |-- gRPC -> Auth Service
  |-- gRPC -> User Service
  `-- gRPC -> Analytics Service
```

Gateway:

- принимает HTTP на `:8080` внутри контейнера;
- в compose опубликован наружу как `localhost:8000`;
- вызывает Auth/User/Analytics по gRPC;
- для protected routes проверяет `Authorization: Bearer <access_token>` через Auth Service;
- передает отчеты Analytics Service клиентам как HTTP JSON.

## HTTP routes

Public:

- `GET /health`
- `POST /api/auth/register`
- `POST /api/auth/verify`
- `POST /api/auth/login`
- `POST /api/auth/refresh`

Protected:

- `POST /api/auth/logout`
- `GET /api/users/:id`
- `PUT /api/users/name`
- `GET /api/analytics/events`
- `GET /api/analytics/reports/registrations`
- `GET /api/analytics/reports/logins`
- `GET /api/analytics/reports/top-users`

Protected routes ожидают access token в заголовке:

```http
Authorization: Bearer <access_token>
```

Refresh route принимает `refresh_token` в JSON body. `access_token` можно передать либо в body, либо через `Authorization`, либо через cookie `access_token`.

## Связь с сервисами

Gateway использует адреса из env:

```env
AUTH_SERVICE_ADDR=auth-service:50051
USER_SERVICE_ADDR=user-service:50052
ANALYTICS_SERVICE_ADDR=analytics-service:50053
```

По умолчанию при локальном `go run` используются:

```env
AUTH_SERVICE_ADDR=localhost:50051
USER_SERVICE_ADDR=localhost:50052
ANALYTICS_SERVICE_ADDR=localhost:50053
```

gRPC proto-файлы лежат в `http_gateway/proto`. Они должны быть синхронизированы с proto контрактами auth, user и analytics сервисов.

## Запуск

Рекомендуемый запуск всего проекта:

```powershell
cd C:\Users\kira4\Loop-company\http_gateway
docker compose up --build
```

После запуска:

- HTTP Gateway: `http://localhost:8000`
- Auth gRPC: `localhost:50051`
- User gRPC: `localhost:50052`
- Analytics gRPC: `localhost:50053`
- Kafka external listener: `localhost:9093`

Локальный запуск только gateway:

```powershell
cd C:\Users\kira4\Loop-company\http_gateway\http_gateway
$env:AUTH_SERVICE_ADDR="localhost:50051"
$env:USER_SERVICE_ADDR="localhost:50052"
$env:ANALYTICS_SERVICE_ADDR="localhost:50053"
go run ./cmd
```

Для такого запуска auth, user и analytics должны уже слушать свои gRPC порты.

## Данные между gateway и сервисами

- Auth flow: HTTP JSON -> Auth gRPC -> HTTP JSON с токенами/статусом.
- User flow: HTTP JSON/path params -> User gRPC -> HTTP JSON profile/settings.
- Analytics flow: HTTP query params -> Analytics gRPC -> HTTP JSON reports/events.
- Kafka events в gateway не проходят: auth/user публикуют события, analytics их читает, gateway только запрашивает отчеты по gRPC.

## Проверки

```powershell
cd C:\Users\kira4\Loop-company\http_gateway\http_gateway
go test ./...
```
