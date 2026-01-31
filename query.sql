-- name: CreateUser :one
INSERT INTO users (google_id, email, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;