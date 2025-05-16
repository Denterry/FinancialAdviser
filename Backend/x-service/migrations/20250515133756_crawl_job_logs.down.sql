-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP INDEX IF EXISTS crawl_job_logs_job_id_idx;
DROP TABLE IF EXISTS crawl_job_logs;
-- +goose StatementEnd