-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS tweets_set_updated_at ON tweets;
DROP TRIGGER IF EXISTS articles_set_updated_at ON articles;
DROP FUNCTION IF EXISTS set_updated_at();
-- +goose StatementEnd