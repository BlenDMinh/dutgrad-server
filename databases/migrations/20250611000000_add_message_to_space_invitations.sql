-- +goose Up
ALTER TABLE space_invitations ADD COLUMN message TEXT;
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
ALTER TABLE space_invitations DROP COLUMN message;
-- +goose StatementBegin
-- +goose StatementEnd
