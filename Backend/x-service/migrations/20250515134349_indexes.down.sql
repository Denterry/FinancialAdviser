-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP INDEX IF EXISTS tweets_created_at_idx;
DROP INDEX IF EXISTS tweets_fulltext_idx;
DROP INDEX IF EXISTS tweet_symbols_symbol_idx;
DROP INDEX IF EXISTS article_symbols_symbol_idx;
-- +goose StatementEnd