-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_query_sessions
ADD COLUMN temp_message TEXT DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_query_sessions
DROP COLUMN temp_message;
-- +goose StatementEnd
