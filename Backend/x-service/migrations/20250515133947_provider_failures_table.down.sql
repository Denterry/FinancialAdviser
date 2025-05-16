-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS provider_failures_trim_trg ON provider_failures;
DROP FUNCTION IF EXISTS provider_failures_trim();
DROP TABLE IF EXISTS provider_failures;
-- +goose StatementEnd