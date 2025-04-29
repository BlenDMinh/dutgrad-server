-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_auth_credentials ALTER COLUMN password_hash DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_auth_credentials ALTER COLUMN password_hash SET NOT NULL;
-- +goose StatementEnd
