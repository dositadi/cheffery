-- +goose Up
BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    deletedAt TIMESTAMPTZ
);

CREATE INDEX idx_name ON users (name text_pattern_ops);

COMMIT;

-- +goose Down
BEGIN;  

DROP INDEX idx_name;
DROP TABLE IF EXISTS users;

COMMIT;
