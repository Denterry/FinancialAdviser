-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP VIEW IF EXISTS v_user_quota;
DROP INDEX IF EXISTS idx_user_usage_user_id;
DROP TABLE IF EXISTS user_usage;
-- +goose StatementEnd
