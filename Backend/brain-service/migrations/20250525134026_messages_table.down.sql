-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_chat_id_created_at;
DROP TABLE IF EXISTS messages;
DROP TYPE IF EXISTS role_enum;
-- +goose StatementEnd
