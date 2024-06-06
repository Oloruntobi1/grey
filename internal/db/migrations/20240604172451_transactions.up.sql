CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    from_user_id UUID NOT NULL,
    to_user_id UUID NOT NULL,
    amount NUMERIC CHECK (amount > 0),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(from_user_id) REFERENCES users(id),
    FOREIGN KEY(to_user_id) REFERENCES users(id)
);