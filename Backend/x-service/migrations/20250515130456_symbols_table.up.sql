-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TYPE symbol_type_enum AS ENUM ('equity', 'crypto', 'etf', 'forex', 'commodity');

CREATE TABLE symbols (
    ticker        TEXT PRIMARY KEY,
    type          symbol_type_enum NOT NULL,
    display_name  TEXT
);

-- COMMENTS
COMMENT ON TABLE symbols IS 'Symbols of tweets';
COMMENT ON COLUMN symbols.ticker IS 'Ticker of the symbol';
COMMENT ON COLUMN symbols.type IS 'Type of the symbol';
COMMENT ON COLUMN symbols.display_name IS 'Display name of the symbol';

-- EXAMPLE DATA
INSERT INTO symbols(ticker,type,display_name) VALUES
 ('AAPL','equity','Apple Inc.'),
 ('TSLA','equity','Tesla'),
 ('BTC','crypto','Bitcoin');
-- +goose StatementEnd