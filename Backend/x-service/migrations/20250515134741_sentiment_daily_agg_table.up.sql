-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TABLE sentiment_daily_agg (
    symbol          TEXT        NOT NULL REFERENCES symbols(ticker),
    day             DATE        NOT NULL,
    avg_score       REAL        NOT NULL,
    pos_cnt         INT         NOT NULL,
    neg_cnt         INT         NOT NULL,
    neu_cnt         INT         NOT NULL,
    PRIMARY KEY (symbol, day)
);

-- COMMENTS
COMMENT ON TABLE sentiment_daily_agg IS 'Sentiment daily aggregate';
COMMENT ON COLUMN sentiment_daily_agg.symbol IS 'Symbol of the sentiment daily aggregate';
COMMENT ON COLUMN sentiment_daily_agg.day IS 'Day of the sentiment daily aggregate';
COMMENT ON COLUMN sentiment_daily_agg.avg_score IS 'Average score of the sentiment daily aggregate';
COMMENT ON COLUMN sentiment_daily_agg.pos_cnt IS 'Positive count of the sentiment daily aggregate';
COMMENT ON COLUMN sentiment_daily_agg.neg_cnt IS 'Negative count of the sentiment daily aggregate';
COMMENT ON COLUMN sentiment_daily_agg.neu_cnt IS 'Neutral count of the sentiment daily aggregate';

-- refresh helper
CREATE OR REPLACE FUNCTION refresh_sentiment_daily_agg() RETURNS void LANGUAGE plpgsql AS $$
BEGIN
    TRUNCATE sentiment_daily_agg;
    INSERT INTO sentiment_daily_agg
    SELECT
        s.symbol,
        date_trunc('day', t.created_at)::date AS day,
        avg(t.sentiment_score)                AS avg_score,
        count(*) FILTER (WHERE t.sentiment_label='POS') AS pos_cnt,
        count(*) FILTER (WHERE t.sentiment_label='NEG') AS neg_cnt,
        count(*) FILTER (WHERE t.sentiment_label='NEU') AS neu_cnt
    FROM symbols s
    JOIN tweet_symbols ts ON ts.symbol = s.ticker
    JOIN tweets t         ON t.id = ts.tweet_id
    WHERE t.sentiment_score IS NOT NULL
    GROUP BY s.symbol, day;
END $$;
-- +goose StatementEnd