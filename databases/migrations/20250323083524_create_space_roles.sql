-- +goose Up
-- +goose StatementBegin
CREATE TABLE space_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    permission INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd
-- +goose Down
DROP TABLE space_roles;
-- +goose StatementBegin
-- +goose StatementEnd
