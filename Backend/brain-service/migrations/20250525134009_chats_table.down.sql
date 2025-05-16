-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chats_user_id;
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd
