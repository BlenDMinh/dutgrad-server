-- +goose Up
-- +goose StatementBegin
ALTER TABLE space_invitations ADD CONSTRAINT uniq_user_space_role UNIQUE(invited_user_id, space_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE space_invitations DROP CONSTRAINT uniq_user_space_role;
-- +goose StatementEnd
