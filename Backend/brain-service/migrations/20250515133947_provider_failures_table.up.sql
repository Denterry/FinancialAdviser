-- +goose Up
-- +migrate Up
CREATE TABLE provider_failures (
    id          bigserial PRIMARY KEY,
    provider    provider_enum NOT NULL,
    at          timestamptz   NOT NULL DEFAULT now(),
    code        int,
    body        text
);

COMMENT ON TABLE provider_failures IS 'Provider failures';
COMMENT ON COLUMN provider_failures.id IS 'Unique identifier for the provider failure';
COMMENT ON COLUMN provider_failures.provider IS 'Provider of the provider failure';
COMMENT ON COLUMN provider_failures.at IS 'Timestamp of the provider failure';
COMMENT ON COLUMN provider_failures.code IS 'Code of the provider failure';
COMMENT ON COLUMN provider_failures.body IS 'Body of the provider failure';

-- retain only last 10 000 rows automatically (ring-buffer)
CREATE OR REPLACE FUNCTION provider_failures_trim() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM provider_failures
    WHERE id < (SELECT max(id) FROM provider_failures) - 10000;
    RETURN NULL;
END $$;

CREATE TRIGGER provider_failures_trim_trg
AFTER INSERT ON provider_failures
EXECUTE PROCEDURE provider_failures_trim();
