# EventBooker Service

**EventBooker** — это высокопроизводительный асинхронный микросервис на Go для управления мероприятиями и бронирования мест с защитой от овербукинга, временно удерживаемыми бронями (TTL) и фоновым воркером автоматической отмены просроченных заявок с использованием **PostgreSQL**.

### Технологический стек:

* **Язык**: Go 1.22+
* **HTTP-Фреймворк**: `ginext` (`wb-go/wbf`)
* **База данных**: PostgreSQL (`wb-go/wbf/dbpg/pgx-driver`)
* **Логирование**: Zap (`wb-go/wbf/logger`)
* **Фронтенд**: HTML5, CSS3, Vanilla JS

## Основные возможности

1. **Создание и управление мероприятиями**: Создание ивентов с указанием общего количества мест и кастомным временем жизни брони (TTL).
2. **Асинхронное фоновое очищение (Worker)**: Воркер в режиме реального времени находит просроченные брони в статусе `pending`, отменяет их и автоматически возвращает места обратно в доступный пул (`available_seats`).
3. **Безопасная обработка бронирования**:
* `pending` — Место временно зарезервировано, ожидает подтверждения пользователя до истечения TTL.
* `confirmed` — Бронь успешно подтверждена.
* `cancelled` — Бронь отменена вручную или автоматически из-за истечения TTL.


4. **Управление пользователями и история**: Регистрация пользователей, просмотр профиля и получение персонального списка бронирований.
5. **Встроенный Dashboard**: Веб-интерфейс для удобного просмотра ивентов, оформления и подтверждения броней.


## Требования к окружению

* **Go**: `1.22` или выше
* **Docker** и **Docker Compose** (для локального запуска PostgreSQL)

---

## Запуск и настройка

### 1. Клонирование репозитория

```bash
git clone https://github.com/QuUteO/EventBooker.git
cd EventBooker
```

### 2. Запуск инфраструктуры (PostgreSQL)

Запустите контейнер с базой данных:

```bash
docker-compose up -d
```

### 3. Применение миграций БД

Убедитесь, что таблицы `users`, `events` и `bookings` созданы в вашей базе данных PostgreSQL:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    telegram_use VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL,
    total_seats INT NOT NULL CHECK (total_seats > 0),
    available_seats INT NOT NULL CHECK (available_seats >= 0),
    booking_ttl_minutes INT NOT NULL DEFAULT 15,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'confirmed', 'cancelled'
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_bookings_status_expires ON bookings(status, expires_at);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_event_id ON bookings(event_id);

-- +goose Down
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;
```

### 4. Конфигурация (`config.yaml`)

Создайте файл `config.yaml` в корне проекта со следующими настройками:

```yaml
env: local

dsn: "postgres://user:user@localhost:5432/eventbooker?sslmode=disable"

http:
  addr: ":8080"
```

### 5. Запуск сервиса

```bash
go run ./cmd/main.go
```

После запуска сервис будет доступен по адресу: `http://localhost:8080`

---

## Веб-интерфейс (UI)

Сервис поставляется со встроенной панелью управления.

Откройте в браузере: **`http://localhost:8080`**

**Функции UI:**

* Просмотр предстоящих мероприятий и доступных мест.
* Быстрое создание пользователя и оформление бронирования.
* Отслеживание таймера истечения срока брони и мгновенное подтверждение/отмена.

---

## API Эндпоинты

### Пользователи (Users)

#### 1. Регистрация пользователя

* **POST** `/api/v1/users`
* **Body**:
```json
{
  "email": "user@example.com",
  "telegram_use": "@username"
}
```



#### 2. Профиль пользователя

* **GET** `/api/v1/users/:id`

#### 3. Все бронирования пользователя

* **GET** `/api/v1/users/:id/bookings`

---

### Мероприятия (Events)

#### 1. Создать мероприятие

* **POST** `/api/v1/events`
* **Body**:
```json
{
  "title": "Go Conference 2026",
  "event_date": "2026-09-15T18:00:00Z",
  "total_seats": 100,
  "booking_ttl_minutes": 15
}
```



#### 2. Список предстоящих мероприятий

* **GET** `/api/v1/events`

#### 3. Детали мероприятия и список броней

* **GET** `/api/v1/events/:id`

---

### Бронирование (Bookings)

#### 1. Забронировать место

* **POST** `/api/v1/events/:id/book`
* **Body**:
```json
{
  "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
}
```


* **Ответ (`201 Created`)**:
```json
{
  "booking_id": "c1f7a012-3456-7890-abcd-ef1234567890",
  "status": "pending",
  "expires_at": "2026-07-30T17:15:00Z"
}
```



#### 2. Подтвердить бронирование

* **POST** `/api/v1/bookings/:id/confirm`

#### 3. Отменить бронирование

* **POST** `/api/v1/bookings/:id/cancel`

---

## Структура проекта

```text
EventBooker/
├── cmd/
│   └── main.go              # Точка входа в приложение (Graceful Shutdown)
├── internal/
│   ├── config/              # Загрузка конфигураций
│   ├── handler/             # HTTP Хэндлеры (ginext / Gin)
│   ├── service/             # Бизнес-логика бронирования
│   ├── repository/          # Взаимодействие с PostgreSQL (pgx)
│   ├── worker/              # Фоновый cleaner просроченных броней
│   └── model/               # DTO, модели и ошибки
├── web/
│   └── index.html           # Веб-интерфейс (Dashboard)
├── logs/                    # Логи работы приложения
├── docker-compose.yml       # Окружение (PostgreSQL)
├── config.yaml              # Конфигурационный файл
├── go.mod
└── README.md
```