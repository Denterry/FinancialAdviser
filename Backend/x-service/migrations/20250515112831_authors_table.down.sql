-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP INDEX  IF EXISTS authors_provider_username_uq;
DROP TABLE  IF EXISTS authors;
DROP TYPE   IF EXISTS provider_enum;
-- +goose StatementEnd
