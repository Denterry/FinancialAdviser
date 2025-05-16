-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tweet_symbols;
-- +goose StatementEnd