-- +goose Up
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE chunks (
    id         TEXT          PRIMARY KEY,
    memory_id  TEXT          NOT NULL REFERENCES memories(id) ON DELETE CASCADE,
    content    TEXT          NOT NULL,
    embedding  vector(1536),
    sequence   INTEGER       NOT NULL,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX idx_chunks_memory_id ON chunks(memory_id);

-- +goose Down
DROP TABLE chunks;
DROP EXTENSION IF EXISTS vector;
