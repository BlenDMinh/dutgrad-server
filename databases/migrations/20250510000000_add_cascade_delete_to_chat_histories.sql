-- +goose Up
-- +goose StatementBegin
-- First, drop the existing foreign key constraint
ALTER TABLE "chat_histories" DROP CONSTRAINT "chat_histories_session_id_fkey";

-- Then, add it back with ON DELETE CASCADE
ALTER TABLE "chat_histories" ADD CONSTRAINT "chat_histories_session_id_fkey"
    FOREIGN KEY ("session_id") REFERENCES "user_query_sessions"("id") ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert to the original constraint without CASCADE
ALTER TABLE "chat_histories" DROP CONSTRAINT "chat_histories_session_id_fkey";
ALTER TABLE "chat_histories" ADD CONSTRAINT "chat_histories_session_id_fkey"
    FOREIGN KEY ("session_id") REFERENCES "user_query_sessions"("id");
-- +goose StatementEnd
