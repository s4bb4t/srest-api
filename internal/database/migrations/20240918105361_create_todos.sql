-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS public.todos (
    id SERIAL PRIMARY KEY,
    title TEXT,
    created TIMESTAMPTZ DEFAULT NOW(),
    is_done BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE IF EXISTS public.todos;
