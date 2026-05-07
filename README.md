# HTTP Gateway

Public HTTP entry point for Loop backend clients.

## Responsibilities

- Accept HTTP requests from web and mobile clients.
- Call auth, user, and analytics services over gRPC.
- Validate protected routes through Auth Service.
- Return analytics reports from Analytics Service.

Gateway does not publish or consume Kafka messages directly.

## Configuration

```env
AUTH_SERVICE_ADDR=auth-service:50051
USER_SERVICE_ADDR=user-service:50052
ANALYTICS_SERVICE_ADDR=analytics-service:50053
```

## Routes

- `POST /api/auth/register`
- `POST /api/auth/verify`
- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `POST /api/auth/logout`
- `GET /api/users/:id`
- `PUT /api/users/name`
- `GET /api/analytics/events`
- `GET /api/analytics/reports/registrations`
- `GET /api/analytics/reports/logins`
- `GET /api/analytics/reports/top-users`

## Local Run

```powershell
docker compose up --build
```
