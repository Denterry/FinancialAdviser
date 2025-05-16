-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd