-- +goose Down
-- +migrate Down
DROP INDEX  IF EXISTS articles_provider_url_uq;
DROP TABLE  IF EXISTS articles;
