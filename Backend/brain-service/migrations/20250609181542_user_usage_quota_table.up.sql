-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_usage (
    user_id         UUID        PRIMARY KEY,
    period_start    DATE        NOT NULL,     -- начало расчётного периода (например, 1-е число месяца)
    tokens_used     BIGINT      NOT NULL DEFAULT 0,
    tokens_limit    BIGINT      NOT NULL DEFAULT 200000  -- дефолтный месячный лимит
);
-- COMMENTS
COMMENT ON TABLE user_usage IS 'Stores usage data for each user';
COMMENT ON COLUMN user_usage.user_id IS 'Unique identifier for the user';
COMMENT ON COLUMN user_usage.period_start IS 'Start date of the usage period';
COMMENT ON COLUMN user_usage.tokens_used IS 'Total tokens consumed by the user during the period';
COMMENT ON COLUMN user_usage.tokens_limit IS 'Maximum allowed tokens for the period';

-- INDEXES
CREATE INDEX IF NOT EXISTS idx_user_usage_user_id ON user_usage(user_id);

-- VIEWS
CREATE OR REPLACE VIEW v_user_quota AS
SELECT
    user_id,
    period_start,
    tokens_used,
    tokens_limit,
    (tokens_limit - tokens_used) AS tokens_left
FROM user_usage;
-- +goose StatementEnd
