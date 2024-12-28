CREATE TYPE user_types AS ENUM('company', 'job_seeker');
CREATE TYPE sexy AS ENUM('f', 'm');

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    email VARCHAR(120) NOT NULL,
    first_name VARCHAR(120),
    last_name VARCHAR(120),
    password VARCHAR(120) NOT NULL,
    is_deleted BOOL DEFAULT false,
    auth_source VARCHAR(120) NOT NULL,
    updated_at TIMESTAMPTZ,
    last_token_update TIMESTAMPTZ,
    verified_email BOOL DEFAULT false,
    user_type user_types DEFAULT 'job_seeker',
    is_banned BOOL DEFAULT false,               -- Указывает, заблокирован ли пользователь.
    date_joined TIMESTAMPTZ DEFAULT NOW(),
    sexy sexy,
    UNIQUE (email)
);

CREATE TABLE phones (
    id UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    user_id UUID UNIQUE NOT NULL,
    number BIGINT UNIQUE NOT NULL,
    country_code VARCHAR(5) NOT NULL,
    verified BOOL DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
