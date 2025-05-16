-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tweets;
-- +goose StatementEnd