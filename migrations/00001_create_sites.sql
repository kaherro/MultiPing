-- SQLBook: Code
-- +goose Up
CREATE TABLE sites (
    id           BIGSERIAL PRIMARY KEY,
    url          TEXT NOT NULL UNIQUE,
    interval_sec INT NOT NULL DEFAULT 60,
    timeout_ms   INT NOT NULL DEFAULT 5000,
    enabled      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE sites;