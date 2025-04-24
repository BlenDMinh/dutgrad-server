-- +goose Up
-- +goose StatementBegin
ALTER TABLE tiers
    RENAME COLUMN doc_per_space_limit TO document_limit;

ALTER TABLE tiers
    ADD file_size_limit_kb INT DEFAULT 5120; -- Maximum file size allowed per document in KB (default: 5MB)

ALTER TABLE tiers
    ADD api_call_limit INT DEFAULT 100; -- Maximum number of API calls allowed per day (default: 100)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tiers
    RENAME COLUMN document_limit TO doc_per_space_limit;

ALTER TABLE tiers
    DROP COLUMN file_size_limit_kb;

ALTER TABLE tiers
    DROP COLUMN api_call_limit;
-- +goose StatementEnd
