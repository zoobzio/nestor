-- +goose Up
CREATE TABLE memories (
    id             TEXT        PRIMARY KEY,
    agent_id       TEXT        NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    content        TEXT        NOT NULL,
    correlation_id TEXT,
    metadata       JSONB,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_memories_agent_id ON memories(agent_id);
CREATE INDEX idx_memories_correlation_id ON memories(correlation_id) WHERE correlation_id IS NOT NULL;

-- +goose Down
DROP TABLE memories;
