-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_llm_logs_message_id;
DROP TABLE IF EXISTS llm_logs;
-- +goose StatementEnd
