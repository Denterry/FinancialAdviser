-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TABLE crawl_job_logs (
    id        BIGSERIAL PRIMARY KEY,
    job_id    BIGINT       REFERENCES crawl_jobs(id) ON DELETE CASCADE,
    level     TEXT         NOT NULL,
    ts        TIMESTAMPTZ  NOT NULL DEFAULT now(),
    message   jsonb        NOT NULL
);

-- COMMENTS
COMMENT ON TABLE crawl_job_logs IS 'Crawl job logs';
COMMENT ON COLUMN crawl_job_logs.id IS 'Unique identifier for the crawl job log';
COMMENT ON COLUMN crawl_job_logs.job_id IS 'ID of the crawl job';
COMMENT ON COLUMN crawl_job_logs.level IS 'Level of the log message (info / warn / error)';
COMMENT ON COLUMN crawl_job_logs.ts IS 'Timestamp of the log message';
COMMENT ON COLUMN crawl_job_logs.message IS 'Message of the log';

-- INDEXES
CREATE INDEX crawl_job_logs_job_id_idx ON crawl_job_logs(job_id);
-- +goose StatementEnd