-- +goose Down
-- +migrate Down
-- +goose StatementBegin
DELETE FROM symbols;
DROP TABLE  IF EXISTS symbols;
DROP TYPE   IF EXISTS symbol_type_enum;
-- +goose StatementEnd
