-- +goose Up
ALTER TABLE user_queries 
ALTER COLUMN query TYPE VARCHAR(1024);
-- +goose StatementBegin
-- +goose StatementEnd
20250526000000
-- +goose Down
ALTER TABLE user_queries 
ALTER COLUMN query TYPE TEXT;
-- +goose StatementBegin
-- +goose StatementEnd
