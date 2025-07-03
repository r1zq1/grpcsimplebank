CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    owner TEXT NOT NULL,
    balance BIGINT NOT NULL CHECK (balance >= 0)
);

CREATE TABLE transfers (
    id BIGSERIAL PRIMARY KEY,
    from_account_id BIGINT NOT NULL REFERENCES accounts(id),
    to_account_id BIGINT NOT NULL REFERENCES accounts(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);