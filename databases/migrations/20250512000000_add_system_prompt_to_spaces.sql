-- Add system_prompt field to spaces table
-- +goose Up
ALTER TABLE spaces ADD COLUMN system_prompt VARCHAR(1024) DEFAULT 'You are an AI assistant for answering questions about documents in this space. Provide helpful, accurate, and concise information based on the content available.';

-- +goose Down
ALTER TABLE spaces DROP COLUMN system_prompt;
