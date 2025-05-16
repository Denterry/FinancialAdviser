-- +goose Up
-- +migrate Up
CREATE INDEX tweets_created_at_idx
    ON tweets(created_at DESC);

CREATE INDEX tweets_fulltext_idx
    ON tweets
    USING GIN (to_tsvector('simple', text));

CREATE INDEX tweet_symbols_symbol_idx
    ON tweet_symbols(symbol);

CREATE INDEX article_symbols_symbol_idx
    ON article_symbols(symbol);
