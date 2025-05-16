CREATE TYPE message_role AS ENUM ('system', 'user', 'assistant');
CREATE TYPE message_status AS ENUM ('pending', 'complete', 'failed');

CREATE TABLE IF NOT EXISTS conversations (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    last_activity BIGINT NOT NULL
);

CREATE INDEX idx_conversations_user_id ON conversations(user_id);
CREATE INDEX idx_conversations_created_at ON conversations(created_at);