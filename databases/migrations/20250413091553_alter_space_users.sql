-- +goose Up
-- +goose StatementBegin
ALTER TABLE space_users ADD CONSTRAINT uniq_user_space UNIQUE(user_id, space_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE space_users DROP CONSTRAINT uniq_user_space;
-- +goose StatementEnd
