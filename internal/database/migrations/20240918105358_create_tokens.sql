-- +goose Up
CREATE TABLE IF NOT EXISTS public.tokens (
    user_id SERIAL PRIMARY KEY, 
    token TEXT NOT NULL,
    date TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS public.tokens;
