-- SQLBook: Code
-- +goose Up
CREATE TABLE check_results (
    id          BIGSERIAL PRIMARY KEY,
    site_id     BIGINT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    ok          BOOLEAN NOT NULL,
    http_code   INT,
    latency_ms  INT,
    error       TEXT,
    checked_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_check_results_site_checked_at
    ON check_results (site_id, checked_at DESC);

-- +goose Down
DROP TABLE check_results;