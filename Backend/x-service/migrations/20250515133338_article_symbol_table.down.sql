-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP TABLE IF EXISTS article_symbols;
-- +goose StatementEnd