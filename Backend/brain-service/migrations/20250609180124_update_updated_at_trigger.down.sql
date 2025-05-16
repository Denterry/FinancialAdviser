-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_chats_updated_at ON chats;
DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd