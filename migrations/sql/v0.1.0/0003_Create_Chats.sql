CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY,
    owner_id UUID REFERENCES user_accounts(id),
    created_at TIMESTAMPTZ NOT NULL,
    title TEXT NOT NULL
);