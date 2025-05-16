-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TYPE role_enum AS ENUM ('user', 'assistant', 'system');

CREATE TABLE messages (
    id          BIGSERIAL PRIMARY KEY,
    chat_id     UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    role        role_enum NOT NULL,
    content     TEXT NOT NULL,
    token_count INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- COMMENTS
COMMENT ON TABLE messages IS 'Messages exchanged in a chat';
COMMENT ON COLUMN messages.id IS 'Message ID';
COMMENT ON COLUMN messages.chat_id IS 'Foreign key to chat';
COMMENT ON COLUMN messages.role IS 'Sender role (user or assistant)';
COMMENT ON COLUMN messages.content IS 'Text content of the message';
COMMENT ON COLUMN messages.token_count IS 'Number of tokens used by the message';
COMMENT ON COLUMN messages.created_at IS 'When the message was sent';

-- INDEXES
CREATE INDEX idx_messages_chat_id_created_at ON messages(chat_id, created_at);
-- +goose StatementEnd
