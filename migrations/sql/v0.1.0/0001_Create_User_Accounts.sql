CREATE TABLE IF NOT EXISTS user_accounts (
    "id" UUID PRIMARY KEY,
    "tag" VARCHAR(32) UNIQUE NOT NULL,
    "name" VARCHAR(32) NOT NULL,
    "desc" TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(15) UNIQUE
);
