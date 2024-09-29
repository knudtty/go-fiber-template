-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE auth_providers AS ENUM ('github', 'google');

-- Create users table
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    email VARCHAR (255) NOT NULL UNIQUE,
    name TEXT NOT NULL,
    password_hash VARCHAR (255) NOT NULL DEFAULT '',
    refresh_token TEXT NOT NULL DEFAULT '',
    user_status INT NOT NULL,
    user_role VARCHAR (25) NOT NULL,
    avatar_url text NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE oauth_accounts (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    provider auth_providers NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL DEFAULT '',
    refresh_token TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (provider, provider_user_id)
);

-- Add index
CREATE INDEX active_users ON users (id) WHERE user_status = 1;
