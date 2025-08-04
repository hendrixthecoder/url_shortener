-- name: CreateURLAnalyticRecord :one
INSERT INTO url_analytics (
    id,
    created_at,
    updated_at,
    referer,
    user_agent,
    ip,
    short_url
) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
