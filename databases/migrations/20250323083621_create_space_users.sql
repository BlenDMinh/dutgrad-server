-- +goose Up
CREATE TABLE space_users (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    space_id INT REFERENCES spaces(id) ON DELETE CASCADE,
    space_role_id INT REFERENCES space_roles(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE space_users;
-- +goose StatementEnd