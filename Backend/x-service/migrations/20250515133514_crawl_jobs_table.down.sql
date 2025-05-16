-- +goose Down
-- +migrate Down
DROP TABLE IF EXISTS crawl_jobs;
DROP TYPE IF EXISTS crawl_status_enum;
