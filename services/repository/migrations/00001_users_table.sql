-- +goose Up
BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash BYTEA NOT NULL,
    createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    deletedAt TIMESTAMPTZ,
    UNIQUE (email)
);

CREATE INDEX idx_name ON users (name text_pattern_ops);

COMMIT;

-- +goose Down
BEGIN;  

DROP INDEX idx_name;
DROP TABLE IF EXISTS users;

COMMIT;
