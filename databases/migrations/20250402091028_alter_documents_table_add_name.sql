-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents ADD COLUMN name VARCHAR(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN name;
-- +goose StatementEnd
