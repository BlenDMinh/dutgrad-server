-- Add description field to documents table
-- +goose Up
ALTER TABLE documents ADD COLUMN description VARCHAR(1024);

-- +goose Down
ALTER TABLE documents DROP COLUMN description;
