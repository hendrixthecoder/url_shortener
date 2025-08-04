-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    plain_url TEXT NOT NULL,
    short_url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_urls_user_id ON urls(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_short_url;
DROP INDEX IF EXISTS idx_user_id;
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
