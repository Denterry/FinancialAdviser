-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP INDEX  IF EXISTS articles_provider_url_uq;
DROP TABLE  IF EXISTS articles;
-- +goose StatementEnd