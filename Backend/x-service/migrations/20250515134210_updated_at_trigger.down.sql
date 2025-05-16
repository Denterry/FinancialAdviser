-- +goose Down
-- +migrate Down
DROP TRIGGER IF EXISTS tweets_set_updated_at ON tweets;
DROP TRIGGER IF EXISTS articles_set_updated_at ON articles;
DROP FUNCTION IF EXISTS set_updated_at();
