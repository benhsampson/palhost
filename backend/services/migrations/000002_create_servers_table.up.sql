CREATE TABLE IF NOT EXISTS servers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    host CIDR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);