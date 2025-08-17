CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY,
    tag VARCHAR(15) NOT NULL UNIQUE,
    owner_id UUID REFERENCES user_accounts(id),
    created_at TIMESTAMPTZ NOT NULL,
    title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_participants (
    chat_id UUID NOT NULL REFERENCES chats(id),
    user_id UUID NOT NULL REFERENCES user_accounts(id),
    PRIMARY KEY (chat_id, user_id)
);