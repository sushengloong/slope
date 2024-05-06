CREATE TABLE messages (
    id VARCHAR(255) PRIMARY KEY,
    body TEXT NOT NULL,
    participant_id VARCHAR(255) NOT NULL,
    participant_type VARCHAR(255) NOT NULL,
    conversation_id VARCHAR(255) NOT NULL,
    metadata JSONB NOT NULL,
    created TIMESTAMPTZ NOT NULL,
    updated TIMESTAMPTZ NOT NULL
);
