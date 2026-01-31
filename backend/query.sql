-- sqlc 用の SQL
-- name: CreateUser :one
INSERT INTO users (google_id, email, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateReport :one
INSERT INTO reports (
    user_id, title, content, learning_notes, is_habit_done, is_public, public_token, submitted_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListReportsByUserID :many
SELECT * FROM reports
WHERE user_id = $1
ORDER BY submitted_date DESC;

-- name: GetReportByPublicToken :one
SELECT * FROM reports
WHERE public_token = $1 AND is_public = true
LIMIT 1;