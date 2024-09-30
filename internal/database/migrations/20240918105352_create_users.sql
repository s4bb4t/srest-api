-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS dev.users (
    id SERIAL PRIMARY KEY, 
    login TEXT UNIQUE,
    username TEXT,
    password TEXT,
    email TEXT UNIQUE,
    date TIMESTAMPTZ DEFAULT NOW(),
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    phone_number TEXT
);

-- +goose Down
DROP TABLE IF EXISTS dev.users;
