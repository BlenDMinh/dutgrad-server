-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents ADD COLUMN mime VARCHAR(255) NOT NULL;
ALTER TABLE documents ADD COLUMN processing_status INT NOT NULL;
ALTER TABLE documents ADD COLUMN size INT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN mime;
ALTER TABLE documents DROP COLUMN processing_status;
ALTER TABLE documents DROP COLUMN size;
-- +goose StatementEnd
