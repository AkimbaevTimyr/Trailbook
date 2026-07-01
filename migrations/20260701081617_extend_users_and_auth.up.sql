-- 4.2.1 Расширение users + auth
ALTER TABLE users
    ADD COLUMN email VARCHAR(255) UNIQUE,
    ADD COLUMN password_hash VARCHAR(255),
    ADD COLUMN avatar_url VARCHAR(500),
    ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Заполнить существующих пользователей временными значениями (для NOT NULL)
UPDATE users SET email = 'user_' || id || '@temp.com', password_hash = '$2a$10$...' WHERE email IS NULL;

ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token       VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
