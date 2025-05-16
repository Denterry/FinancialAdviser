-- +goose Up
-- +migrate Up
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END $$;

CREATE TRIGGER tweets_set_updated_at
BEFORE UPDATE ON tweets
FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER articles_set_updated_at
BEFORE UPDATE ON articles
FOR EACH ROW EXECUTE PROCEDURE set_updated_at();