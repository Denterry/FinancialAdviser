-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS refresh_sentiment_daily_agg();
DROP TABLE IF EXISTS sentiment_daily_agg;
-- +goose StatementEnd