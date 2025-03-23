-- +goose Up
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    space_id INT REFERENCES spaces(id) ON DELETE CASCADE,
    s3_url TEXT NOT NULL,
    privacy_status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE documents;
-- +goose StatementBegin
-- +goose StatementEnd
