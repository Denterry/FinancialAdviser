-- +goose Up
-- +migrate Up
CREATE TABLE tweet_symbols (
    tweet_id UUID REFERENCES tweets(id) ON DELETE CASCADE,
    symbol   TEXT   NOT NULL REFERENCES symbols(ticker) ON DELETE CASCADE,
    PRIMARY KEY (tweet_id, symbol)
);

COMMENT ON TABLE tweet_symbols IS 'Tweet symbols';
COMMENT ON COLUMN tweet_symbols.tweet_id IS 'Tweet ID';
COMMENT ON COLUMN tweet_symbols.symbol IS 'Symbol';
