-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE action_enum AS ENUM ('entrance', 'exit');

CREATE TABLE IF NOT EXISTS actions
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid       UUID NOT NULL,
    action     action_enum,
    created_at TIMESTAMP        DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
DROP TYPE IF EXISTS action_enum;
DROP EXTENSION IF EXISTS pgcrypto;
