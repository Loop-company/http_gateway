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
- в локальной инфраструктуре опубликован наружу как `localhost:8000`;
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

## Связь с сервисами

Gateway использует адреса из env:

```env
AUTH_SERVICE_ADDR=auth-service:50051
USER_SERVICE_ADDR=user-service:50052
ANALYTICS_SERVICE_ADDR=analytics-service:50053
```

gRPC proto-файлы лежат в `http_gateway/proto`. Они должны быть синхронизированы с proto контрактами auth, user и analytics сервисов.

## Запуск

Docker Compose для локального запуска вынесен в репозиторий `loop_infra`.

Из корня `loop_infra`:

```bash
docker compose up --build
```

В этом репозитории остается код сервиса, `Dockerfile` и пример переменных окружения `http_gateway/.env.example`.
