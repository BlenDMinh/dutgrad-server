-- +goose Up
CREATE TABLE user_queries (
    id SERIAL PRIMARY KEY,
    query_session_id INT REFERENCES user_query_sessions(id) ON DELETE CASCADE,
    query TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
-- +goose StatementBegin
-- +goose StatementEnd
20250323074259
-- +goose Down
DROP TABLE user_queries;
-- +goose StatementBegin
-- +goose StatementEnd