-- Update description field length in documents table
-- +goose Up
ALTER TABLE documents ALTER COLUMN description TYPE VARCHAR(8192);

-- +goose Down
ALTER TABLE documents ALTER COLUMN description TYPE VARCHAR(1024);
