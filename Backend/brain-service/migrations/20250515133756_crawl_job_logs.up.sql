-- +goose Up
-- +migrate Up
CREATE TABLE crawl_job_logs (
    id        bigserial PRIMARY KEY,
    job_id    bigint       REFERENCES crawl_jobs(id) ON DELETE CASCADE,
    level     text         NOT NULL,          -- info / warn / error
    ts        timestamptz  NOT NULL DEFAULT now(),
    message   jsonb        NOT NULL
);

COMMENT ON TABLE crawl_job_logs IS 'Crawl job logs';
COMMENT ON COLUMN crawl_job_logs.id IS 'Unique identifier for the crawl job log';
COMMENT ON COLUMN crawl_job_logs.job_id IS 'ID of the crawl job';
COMMENT ON COLUMN crawl_job_logs.level IS 'Level of the log message';
COMMENT ON COLUMN crawl_job_logs.ts IS 'Timestamp of the log message';
COMMENT ON COLUMN crawl_job_logs.message IS 'Message of the log';

CREATE INDEX crawl_job_logs_job_id_idx ON crawl_job_logs(job_id);
