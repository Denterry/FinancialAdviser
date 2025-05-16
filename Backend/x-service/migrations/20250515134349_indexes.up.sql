-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE INDEX idx_tweets_created_at
    ON tweets(created_at DESC);

CREATE INDEX idx_tweets_fulltext
    ON tweets
    USING GIN (to_tsvector('simple', text));

CREATE INDEX idx_tweet_symbols_symbol
    ON tweet_symbols(symbol);

CREATE INDEX idx_article_symbols_symbol
    ON article_symbols(symbol);
-- +goose StatementEnd