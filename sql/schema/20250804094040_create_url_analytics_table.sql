-- +goose Up
-- +goose StatementBegin
CREATE TABLE url_analytics (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    referer TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    ip TEXT NOT NULL,
    short_url TEXT NOT NULL REFERENCES urls(short_url) ON DELETE CASCADE
);

CREATE INDEX idx_url_short ON url_analytics(short_url)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_url_short;
DROP TABLE IF EXISTS url_analytics;
-- +goose StatementEnd
