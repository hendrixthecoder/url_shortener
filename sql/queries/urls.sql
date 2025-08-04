-- name: CreateNewShortURL :one
INSERT INTO urls (id, created_at, updated_at, user_id, short_url, plain_url) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetURLEntryByShortURL :one
SELECT * FROM urls WHERE short_url = $1;

-- name: GetURLEntriesByUserID :many
SELECT * FROM urls WHERE user_id = $1;