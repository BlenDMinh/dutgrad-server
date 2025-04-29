-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_users_name_trgm ON users USING gin (username gin_trgm_ops);
CREATE INDEX idx_users_email_trgm ON users USING gin (email gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_name_trgm;
DROP INDEX IF EXISTS idx_users_email_trgm;
-- +goose StatementEnd
