-- +goose Up
CREATE TABLE space_invitations (
    id SERIAL PRIMARY KEY,
    space_id INT REFERENCES spaces(id) ON DELETE CASCADE,
    space_role_id INT REFERENCES space_roles(id) ON DELETE SET NULL,
    invited_user_id INT REFERENCES users(id) ON DELETE CASCADE,
    inviter_id INT REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE space_invitations;
-- +goose StatementBegin
-- +goose StatementEnd
