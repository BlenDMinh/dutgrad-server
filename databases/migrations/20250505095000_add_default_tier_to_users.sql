-- +goose Up
ALTER TABLE users ALTER COLUMN tier_id SET DEFAULT 1;
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
ALTER TABLE users ALTER COLUMN tier_id DROP DEFAULT;
-- +goose StatementBegin
-- +goose StatementEnd