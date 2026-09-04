-- name: CreateUser :one
INSERT INTO users (
    name, email, password_hash
) VALUES (
    $1, $2, $3
)
RETURNING id, createdAt;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 AND deletedAt IS NULL;

-- name: UpdateUser :exec
UPDATE users
SET name = $1, email = $2, password_hash = $3, version = version + 1, updated_at = NOW()
WHERE id = $4 AND version = $5;