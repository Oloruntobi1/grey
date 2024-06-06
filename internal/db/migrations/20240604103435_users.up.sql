CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE
);