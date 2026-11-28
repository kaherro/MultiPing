-- SQLBook: Code
-- +goose Up
CREATE TABLE jobs (
    id           BIGSERIAL PRIMARY KEY,
    site_id      BIGINT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    agent_id     BIGINT NOT NULL REFERENCES agents(id),
    status       TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'assigned', 'completed', 'failed')),

    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_at  TIMESTAMPTZ,
    deadline_at  TIMESTAMPTZ,
    finished_at  TIMESTAMPTZ
);

CREATE INDEX idx_jobs_pending
    ON jobs (created_at, id)
    WHERE status = 'pending';

CREATE INDEX idx_jobs_agent_status
    ON jobs (agent_id, status);

-- +goose Down
DROP TABLE IF EXISTS jobs;