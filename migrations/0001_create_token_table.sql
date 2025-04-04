-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE status_enum AS ENUM ('pending', 'scanned', 'expired');

CREATE TABLE IF NOT EXISTS tokens
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token      TEXT        NOT NULL,
    uuid       UUID        NOT NULL,
    status     status_enum NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    scanned_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
DROP TYPE IF EXISTS status_enum;
DROP EXTENSION IF EXISTS pgcrypto;
