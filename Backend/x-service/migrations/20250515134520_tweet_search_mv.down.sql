-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS refresh_tweet_search_mv();
DROP MATERIALIZED VIEW IF EXISTS tweet_search_mv;
-- +goose StatementEnd
