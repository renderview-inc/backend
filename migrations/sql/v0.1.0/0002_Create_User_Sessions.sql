CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id),
    refresh_token_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    refresh_expires_at TIMESTAMPTZ NOT NULL,
    last_used_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL,
    rotated_from_session_id UUID REFERENCES user_sessions(id)
);

CREATE TABLE IF NOT EXISTS user_login_histories (
    login_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts("id"),
    login_time TIMESTAMPTZ NOT NULL,
    user_agent TEXT NOT NULL,
    ip_address INET NOT NULL,
    success BOOL NOT NULL 
);
