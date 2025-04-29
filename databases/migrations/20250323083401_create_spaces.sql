-- +goose Up
CREATE TABLE spaces (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    privacy_status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE spaces;
-- +goose StatementBegin
-- +goose StatementEnd
