-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TABLE articles (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider      provider_enum NOT NULL,
    url           TEXT          NOT NULL,
    title         TEXT,
    html          TEXT,
    author        TEXT,
    published_at  TIMESTAMPTZ,
    fetched_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),

    sentiment_score real,
    sentiment_label char(3),

    raw_payload   jsonb
);

-- COMMENTS
COMMENT ON TABLE articles IS 'Articles';
COMMENT ON COLUMN articles.id IS 'Unique identifier for the article';
COMMENT ON COLUMN articles.provider IS 'Provider of the article';
COMMENT ON COLUMN articles.url IS 'URL of the article';
COMMENT ON COLUMN articles.title IS 'Title of the article';
COMMENT ON COLUMN articles.html IS 'HTML content of the article';
COMMENT ON COLUMN articles.author IS 'Author of the article';
COMMENT ON COLUMN articles.published_at IS 'Published at timestamp';
COMMENT ON COLUMN articles.fetched_at IS 'Fetched at timestamp';
COMMENT ON COLUMN articles.updated_at IS 'Updated at timestamp';
COMMENT ON COLUMN articles.sentiment_score IS 'Sentiment score of the article';
COMMENT ON COLUMN articles.sentiment_label IS 'Sentiment label of the article';
COMMENT ON COLUMN articles.raw_payload IS 'Raw payload of the article';

-- INDEXES
CREATE UNIQUE INDEX articles_provider_url_uq
    ON articles(provider, url);
-- +goose StatementEnd