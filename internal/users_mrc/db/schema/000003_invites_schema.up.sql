CREATE TABLE invites (
    id UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    invite_code VARCHAR(255) UNIQUE,
    is_used BOOL DEFAULT FALSE,
    group_id INTEGER,
    used_by_user_id UUID UNIQUE,
    created_by_user_id UUID,
    expiration_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (used_by_user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL
);