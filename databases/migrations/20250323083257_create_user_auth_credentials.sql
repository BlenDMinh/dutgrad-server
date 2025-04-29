-- +goose Up
CREATE TABLE "user_auth_credentials" (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    auth_type VARCHAR(50) NOT NULL,
    password_hash TEXT NOT NULL,
    external_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE user_auth_credentials;
-- +goose StatementBegin
-- +goose StatementEnd