CREATE TABLE conversations (
    id VARCHAR(255) PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL,
    channel VARCHAR(255) NOT NULL,
    metadata JSONB NOT NULL,
    status VARCHAR(255) NOT NULL,
    created TIMESTAMPTZ NOT NULL,
    updated TIMESTAMPTZ NOT NULL
);
