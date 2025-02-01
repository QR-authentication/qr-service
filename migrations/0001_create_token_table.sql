-- +goose Up
CREATE TYPE status_enum AS ENUM ('pending', 'scanned', 'expired');

CREATE TABLE IF NOT EXISTS tokens
(
    id         UUID PRIMARY KEY,
    token      TEXT        NOT NULL,
    uuid       UUID        NOT NULL,
    status     status_enum NOT NULL DEFAULT 'pending',
    ip_address VARCHAR(45) NOT NULL,
    scanned_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
DROP TYPE IF EXISTS status_enum;
