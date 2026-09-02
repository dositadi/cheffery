-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password BYTEA NOT NULL,
    createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deletedAt TIMESTAMPTZ,
    UNIQUE (email)
);

CREATE IF NOT EXISTS INDEX idx_name ON users (name text_pattern_ops);

-- +goose Down
DROP TABLE IF EXISTS users;
