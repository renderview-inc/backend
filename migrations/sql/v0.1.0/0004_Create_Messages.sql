CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES user_accounts(id),
    chat_tag VARCHAR(15) REFERENCES chats(tag),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);