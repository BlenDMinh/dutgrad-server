-- +goose Up
CREATE TABLE space_api_keys (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    space_id INT REFERENCES spaces (id) ON DELETE CASCADE
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE space_api_keys;
-- +goose StatementEnd
