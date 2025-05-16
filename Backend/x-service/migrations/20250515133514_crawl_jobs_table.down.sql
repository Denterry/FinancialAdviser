-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP TABLE IF EXISTS crawl_jobs;
DROP TYPE IF EXISTS crawl_status_enum;
-- +goose StatementEnd