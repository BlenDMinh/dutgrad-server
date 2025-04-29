-- +goose Up
CREATE TABLE tiers (
    id SERIAL PRIMARY KEY,
    space_limit INT NOT NULL,
    doc_per_space_limit INT NOT NULL,
    query_history_limit INT NOT NULL,
    query_limit INT NOT NULL,
    cost_month DECIMAL(10,2) NOT NULL,
    discount DECIMAL(5,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE tiers;
-- +goose StatementBegin
-- +goose StatementEnd