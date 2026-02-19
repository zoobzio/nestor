-- +goose Up
CREATE TABLE ingest_jobs (
    id             TEXT        PRIMARY KEY,
    customer_id    TEXT        NOT NULL,
    agent_id       TEXT        NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    memory_id      TEXT        NOT NULL,
    content        TEXT        NOT NULL,
    stage          TEXT        NOT NULL DEFAULT 'validate',
    status         TEXT        NOT NULL DEFAULT 'pending',
    chunk_ids      TEXT[],
    chunk_count    INTEGER     NOT NULL DEFAULT 0,
    embedded_count INTEGER     NOT NULL DEFAULT 0,
    error          TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at     TIMESTAMPTZ,
    completed_at   TIMESTAMPTZ
);

CREATE INDEX idx_ingest_jobs_status ON ingest_jobs(status);
CREATE INDEX idx_ingest_jobs_agent_id ON ingest_jobs(agent_id);

-- +goose Down
DROP TABLE ingest_jobs;
