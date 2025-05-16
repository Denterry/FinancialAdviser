-- +goose Up
-- +migrate Up
-- +goose StatementBegin
CREATE TYPE crawl_status_enum AS ENUM ('running','success','failed');

CREATE TABLE crawl_jobs (
    id               BIGSERIAL PRIMARY KEY,
    provider         provider_enum NOT NULL,
    status           crawl_status_enum NOT NULL DEFAULT 'running',
    started_at       TIMESTAMPTZ   NOT NULL DEFAULT now(),
    finished_at      TIMESTAMPTZ,
    rows_ingested    INT,
    error_text       TEXT
);

-- COMMENTS
COMMENT ON TABLE crawl_jobs IS 'Crawl jobs';
COMMENT ON COLUMN crawl_jobs.id IS 'Unique identifier for the crawl job';
COMMENT ON COLUMN crawl_jobs.provider IS 'Provider of the crawl job';
COMMENT ON COLUMN crawl_jobs.status IS 'Status of the crawl job';
COMMENT ON COLUMN crawl_jobs.started_at IS 'Started at timestamp';
COMMENT ON COLUMN crawl_jobs.finished_at IS 'Finished at timestamp';
COMMENT ON COLUMN crawl_jobs.rows_ingested IS 'Number of rows ingested';
COMMENT ON COLUMN crawl_jobs.error_text IS 'Error text';
-- +goose StatementEnd