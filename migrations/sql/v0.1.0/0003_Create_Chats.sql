CREATE TABLE IF NOT EXISTS chats (
    tag VARCHAR(15) PRIMARY KEY,
    owner_id UUID REFERENCES user_accounts(id),
    created_at TIMESTAMPTZ NOT NULL,
    title TEXT NOT NULL
);