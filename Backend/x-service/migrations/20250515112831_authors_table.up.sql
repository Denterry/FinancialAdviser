-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TYPE provider_enum AS ENUM ('twitter', 'tradingview', 'reddit', 'rss');

CREATE TABLE authors (
    id            VARCHAR(64) PRIMARY KEY,
    username      TEXT        NOT NULL,
    display_name  TEXT,
    verified      BOOLEAN     NOT NULL DEFAULT false,
    provider      provider_enum NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- COMMENTS
COMMENT ON TABLE authors IS 'Authors of tweets';
COMMENT ON COLUMN authors.id IS 'Unique identifier for the author';
COMMENT ON COLUMN authors.username IS 'Username of the author';
COMMENT ON COLUMN authors.display_name IS 'Display name of the author';
COMMENT ON COLUMN authors.verified IS 'Whether the author is verified';
COMMENT ON COLUMN authors.provider IS 'Provider of the author';
COMMENT ON COLUMN authors.created_at IS 'Timestamp when the author was created';

-- INDEXES
CREATE UNIQUE INDEX authors_provider_username_uq
    ON authors (provider, lower(username));
-- +goose StatementEnd