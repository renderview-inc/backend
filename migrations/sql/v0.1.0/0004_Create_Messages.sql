CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    reply_to UUID REFERENCES messages,
    user_id UUID REFERENCES user_accounts(id),
    chat_tag VARCHAR(15) REFERENCES chats(tag),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (reply_to <> id)
);
