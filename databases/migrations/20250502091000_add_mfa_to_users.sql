-- +goose Up
-- Add MFA enabled flag to users table
ALTER TABLE users ADD COLUMN mfa_enabled BOOLEAN DEFAULT FALSE;

-- Create table for storing MFA secrets and recovery codes
CREATE TABLE user_mfas (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    secret VARCHAR(255) NOT NULL,
    backup_codes JSONB DEFAULT '[]',
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_mfa;
ALTER TABLE users DROP COLUMN mfa_enabled;
-- +goose StatementEnd
