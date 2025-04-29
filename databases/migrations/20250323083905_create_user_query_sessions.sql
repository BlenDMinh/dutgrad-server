-- +goose Up
CREATE TABLE user_query_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    space_id INT REFERENCES spaces(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE user_query_sessions;
-- +goose StatementBegin
-- +goose StatementEnd
