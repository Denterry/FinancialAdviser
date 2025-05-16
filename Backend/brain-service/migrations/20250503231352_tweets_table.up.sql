-- +goose Up
-- +migrate Up
CREATE TABLE IF NOT EXISTS tweets (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id   VARCHAR(64)  NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    username    TEXT         NOT NULL,
    text        TEXT         NOT NULL,
    lang        VARCHAR(8),                              -- ISO-639-1 or empty

    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),     -- when tweet was posted
    fetched_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),     -- when we ingested it
    updated_at  TIMESTAMPTZ NOT NULL,

    likes       INT          NOT NULL DEFAULT 0 CHECK (likes    >= 0),
    replies     INT          NOT NULL DEFAULT 0 CHECK (replies  >= 0),
    retweets    INT          NOT NULL DEFAULT 0 CHECK (retweets >= 0),
    views       INT          NOT NULL DEFAULT 0 CHECK (views    >= 0),

    urls        TEXT[]       NOT NULL DEFAULT '{}',
    photos      TEXT[]       NOT NULL DEFAULT '{}',
    videos      TEXT[]       NOT NULL DEFAULT '{}',

    is_financial    BOOLEAN  NOT NULL DEFAULT true,      -- if tweet contains financial data
    sentiment_score REAL,
    sentiment_label CHAR(3),

    raw_json    JSONB                                    -- original payload
);

COMMENT ON TABLE  tweets                   IS 'Ingested tweets / X posts';
COMMENT ON COLUMN tweets.id                IS 'Snowflake Tweet ID';
COMMENT ON COLUMN tweets.author_id         IS 'Snowflake User ID of author';
COMMENT ON COLUMN tweets.username          IS 'Author handle (without @)';
COMMENT ON COLUMN tweets.text              IS 'Full tweet text';
COMMENT ON COLUMN tweets.lang              IS 'Detected or provided language tag (ISO-639-1)';
COMMENT ON COLUMN tweets.created_at        IS 'Timestamp when tweet was originally published';
COMMENT ON COLUMN tweets.fetched_at        IS 'Timestamp when our service stored this tweet';
COMMENT ON COLUMN tweets.likes             IS 'Number of likes at fetch time';
COMMENT ON COLUMN tweets.replies           IS 'Number of replies at fetch time';
COMMENT ON COLUMN tweets.retweets          IS 'Number of retweets/reposts at fetch time';
COMMENT ON COLUMN tweets.views             IS 'Public view count at fetch time';
COMMENT ON COLUMN tweets.urls              IS 'Array of URLs contained in tweet text';
COMMENT ON COLUMN tweets.photos            IS 'Array of photo URLs (media)';
COMMENT ON COLUMN tweets.videos            IS 'Array of video URLs (media)';
COMMENT ON COLUMN tweets.raw_json          IS 'Original JSON payload from twitter-scraper or official API';
COMMENT ON COLUMN tweets.is_financial      IS 'Whether the tweet contains financial data';
COMMENT ON COLUMN tweets.sentiment_score   IS 'Sentiment score of the tweet';
COMMENT ON COLUMN tweets.sentiment_label   IS 'Sentiment label of the tweet';
