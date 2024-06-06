CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    user_id UUID NOT NULL,
    balance NUMERIC CHECK (balance >= 0),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(user_id) REFERENCES users(id)
);