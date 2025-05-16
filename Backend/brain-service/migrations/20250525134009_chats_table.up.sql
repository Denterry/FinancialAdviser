-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TABLE chats (
    id          SERIAL PRIMARY KEY,
    user_id     UUID NOT NULL,
    title       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- COMMENTS
COMMENT ON TABLE chats IS 'Stores user chat sessions';
COMMENT ON COLUMN chats.id IS 'Auto-incremented chat session ID';
COMMENT ON COLUMN chats.user_id IS 'Reference to external user ID';
COMMENT ON COLUMN chats.title IS 'Optional title of the chat';
COMMENT ON COLUMN chats.created_at IS 'Chat creation timestamp';
COMMENT ON COLUMN chats.updated_at IS 'Last update timestamp';

-- INDEXES
CREATE INDEX idx_chats_user_id ON chats(user_id);
-- +goose StatementEnd
