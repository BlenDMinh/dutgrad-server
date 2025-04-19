-- +goose Up
-- +goose StatementBegin
CREATE TABLE "chat_histories" (
    "id" SERIAL PRIMARY KEY,
    "session_id" INT REFERENCES user_query_sessions(id),
    "message" JSONB NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger function for updating updated_at column
CREATE OR REPLACE FUNCTION update_chat_histories_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger that will call the function before update
CREATE TRIGGER update_chat_histories_updated_at
BEFORE UPDATE ON "chat_histories"
FOR EACH ROW
EXECUTE FUNCTION update_chat_histories_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_chat_histories_updated_at ON "chat_histories";
DROP FUNCTION IF EXISTS update_chat_histories_updated_at();
DROP TABLE "chat_histories";
-- +goose StatementEnd
