# HTTP Gateway

HTTP Gateway - публичная точка входа для Web/iOS/Android клиентов. Он принимает HTTP-запросы, валидирует protected routes через Auth Service и обращается к внутренним сервисам по gRPC. Kafka напрямую gateway не использует.

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

Refresh route принимает `refresh_token` в JSON body. `access_token` можно передать в body, через `Authorization` или через cookie `access_token`.

## Связь с сервисами

Gateway использует адреса из env:

```env
AUTH_SERVICE_ADDR=auth-service:50051
USER_SERVICE_ADDR=user-service:50052
ANALYTICS_SERVICE_ADDR=analytics-service:50053
```

gRPC proto-файлы лежат в `http_gateway/proto`. Они должны быть синхронизированы с proto контрактами auth, user и analytics сервисов.

## Запуск

Для запуска всего backend-стека через Docker Compose из корня репозитория:

```bash
docker compose up --build
```

Compose поднимает HTTP Gateway, Auth Service, User Service, Analytics Service, PostgreSQL, Redis, Kafka и нужные Kafka topics.

После запуска:

- HTTP Gateway: `http://localhost:8000`
- Auth gRPC: `localhost:50051`
- User gRPC: `localhost:50052`
- Analytics gRPC: `localhost:50053`
- Kafka external listener: `localhost:9093`

## Данные между gateway и сервисами

- Auth flow: HTTP JSON -> Auth gRPC -> HTTP JSON с токенами/статусом.
- User flow: HTTP JSON/path params -> User gRPC -> HTTP JSON profile/settings.
- Analytics flow: HTTP query params -> Analytics gRPC -> HTTP JSON reports/events.
- Kafka events в gateway не проходят: auth/user публикуют события, analytics их читает, gateway только запрашивает отчеты по gRPC.
