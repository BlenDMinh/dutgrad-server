-- +goose Up
-- +goose StatementBegin
ALTER TABLE spaces 
    ADD document_limit INT DEFAULT 10, -- Maximum number of documents allowed in this space (default: 10)
    ADD file_size_limit_kb INT DEFAULT 5120, -- Maximum file size allowed per document in KB (default: 5MB)
    ADD api_call_limit INT DEFAULT 100; -- Maximum number of API calls allowed per day (default: 100)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE spaces 
    DROP COLUMN document_limit,
    DROP COLUMN file_size_limit_kb,
    DROP COLUMN api_call_limit;
-- +goose StatementEnd
