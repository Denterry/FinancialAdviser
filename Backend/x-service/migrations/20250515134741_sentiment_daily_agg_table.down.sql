-- +goose Down
-- +migrate Down
DROP FUNCTION IF EXISTS refresh_sentiment_daily_agg();
DROP TABLE IF EXISTS sentiment_daily_agg;
