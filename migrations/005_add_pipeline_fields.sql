-- +goose Up
-- Add status field to memories for pipeline tracking
ALTER TABLE memories ADD COLUMN status TEXT NOT NULL DEFAULT 'processing';
CREATE INDEX idx_memories_status ON memories(status);

-- Add chisel metadata fields to chunks
ALTER TABLE chunks ADD COLUMN kind TEXT NOT NULL DEFAULT 'section';
ALTER TABLE chunks ADD COLUMN symbol TEXT;
ALTER TABLE chunks ADD COLUMN context TEXT[];
ALTER TABLE chunks ADD COLUMN start_line INTEGER NOT NULL DEFAULT 0;
ALTER TABLE chunks ADD COLUMN end_line INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE chunks DROP COLUMN end_line;
ALTER TABLE chunks DROP COLUMN start_line;
ALTER TABLE chunks DROP COLUMN context;
ALTER TABLE chunks DROP COLUMN symbol;
ALTER TABLE chunks DROP COLUMN kind;

DROP INDEX idx_memories_status;
ALTER TABLE memories DROP COLUMN status;
