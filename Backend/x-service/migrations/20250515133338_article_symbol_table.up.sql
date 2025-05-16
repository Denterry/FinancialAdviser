-- +goose Up
-- +migrate Up
CREATE TABLE article_symbols (
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    symbol     TEXT NOT NULL REFERENCES symbols(ticker) ON DELETE CASCADE,
    PRIMARY KEY (article_id, symbol)
);

COMMENT ON TABLE article_symbols IS 'Article symbols';
COMMENT ON COLUMN article_symbols.article_id IS 'Article ID';
COMMENT ON COLUMN article_symbols.symbol IS 'Symbol';
