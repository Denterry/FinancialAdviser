-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS llm_logs (
    id                  BIGSERIAL PRIMARY KEY,
    message_id          BIGINT          NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    provider            TEXT            NOT NULL, -- 'openai'
    model               TEXT            NOT NULL, -- 'gpt-4o-mini'
    prompt_tokens       INTEGER         NOT NULL,
    completion_tokens   INTEGER         NOT NULL,
    temperature         NUMERIC(4,2)    NOT NULL,
    top_p               NUMERIC(4,2)    NOT NULL DEFAULT 1.00,
    latency_ms          INTEGER         NOT NULL,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- COMMENTS
COMMENT ON TABLE llm_logs IS 'Stores logs of llms used for messages';
COMMENT ON COLUMN llm_logs.id IS 'Auto-incrementing primary key identifier';
COMMENT ON COLUMN llm_logs.message_id IS 'Foreign key reference to messages table (ON DELETE CASCADE)';
COMMENT ON COLUMN llm_logs.provider IS 'LLM provider name (e.g. openai, anthropic, mistral)';
COMMENT ON COLUMN llm_logs.model IS 'Specific model version used (e.g. gpt-4-turbo, claude-3-opus)';
COMMENT ON COLUMN llm_logs.prompt_tokens IS 'Number of tokens consumed by the prompt';
COMMENT ON COLUMN llm_logs.completion_tokens IS 'Number of tokens generated in the response';
COMMENT ON COLUMN llm_logs.temperature IS 'Sampling temperature (0.00-2.00) controlling randomness';
COMMENT ON COLUMN llm_logs.top_p IS 'Nucleus sampling threshold (0.00-1.00) for response diversity';
COMMENT ON COLUMN llm_logs.latency_ms IS 'Total request latency in milliseconds';
COMMENT ON COLUMN llm_logs.created_at IS 'Timestamp of log entry creation (UTC timezone)';

-- INDEXES
CREATE INDEX IF NOT EXISTS idx_llm_logs_message_id ON llm_logs(message_id);
-- +goose StatementEnd
