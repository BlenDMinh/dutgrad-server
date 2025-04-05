-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN mime;
ALTER TABLE documents ADD COLUMN mime_type VARCHAR(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN mime_type;
ALTER TABLE documents ADD COLUMN mime VARCHAR(255) NOT NULL;
-- +goose StatementEnd
