-- +goose Up
CREATE TABLE filter_groups (
    id              TEXT        PRIMARY KEY,
    agent_id        TEXT        NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    name            TEXT        NOT NULL,
    criteria        JSONB,
    prompt_template TEXT        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_filter_groups_agent_id ON filter_groups(agent_id);

-- +goose Down
DROP TABLE filter_groups;
