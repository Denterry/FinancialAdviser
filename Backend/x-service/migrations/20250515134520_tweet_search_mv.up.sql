-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW tweet_search_mv AS
SELECT  t.id,
        t.text,
        t.lang,
        t.created_at,
        t.sentiment_score,
        t.sentiment_label,
        a.username,
        array_agg(ts.symbol) AS symbols
FROM tweets t
JOIN authors a          ON a.id = t.author_id
LEFT JOIN tweet_symbols ts ON ts.tweet_id = t.id
GROUP BY t.id, a.username;

CREATE INDEX tweet_search_mv_created_at_idx
    ON tweet_search_mv(created_at DESC);

-- refresh helper
CREATE OR REPLACE FUNCTION refresh_tweet_search_mv() RETURNS void LANGUAGE sql AS
$$ REFRESH MATERIALIZED VIEW CONCURRENTLY tweet_search_mv; $$;
-- +goose StatementEnd