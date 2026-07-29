-- +goose Up
-- 1. Таблица пользователей
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    telegram_use VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Таблица мероприятий
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL,
    total_seats INT NOT NULL CHECK (total_seats > 0),
    available_seats INT NOT NULL CHECK (available_seats >= 0),
    booking_ttl_minutes INT NOT NULL DEFAULT 15,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Таблица бронирований
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, confirmed, canceled
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Для метода ListBookingsByEventID
CREATE INDEX idx_bookings_event_id ON bookings (event_id);

-- Для метода CancelExpiredBookings
CREATE INDEX idx_bookings_expires_at_pending
    ON bookings (expires_at)
    WHERE status = 'pending';

CREATE INDEX idx_bookings_user_id ON bookings (user_id);

-- +goose Down
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;