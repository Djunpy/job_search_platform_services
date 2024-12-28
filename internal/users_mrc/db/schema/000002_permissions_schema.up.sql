CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    codename VARCHAR(255)
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE user_groups (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    group_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE TABLE group_permissions (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE TABLE user_permissions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    permission_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, permission_id)
);

INSERT INTO permissions(name, codename)
VALUES
    ('Add Post', 'add_post'),
    ('Change Post', 'change_post'),
    ('Delete Post', 'delete_post'),
    ('Hide Post', 'hide_post'),
    ('Block Post', 'block_post'),
    ('View Post', 'view_post');

INSERT INTO groups (name)
VALUES
    ('administrators'),
    ('moderators'),
    ('ordinary_users'),
    ('contributors'),
    ('product_managers'),
    ('project_managers'),
    ('premium_users'),
    ('testers '),
    ('supports'),
    ('developers'),
    ('companies'),
    ('job_seekers');

-- INSERT INTO group_permissions (group_id, permission_id)
-- SELECT
--     (SELECT id FROM groups WHERE name = 'administrators'),
--     id
-- FROM permissions
-- WHERE codename IN ('add_post', 'edit_post', 'delete_post', 'view_post', 'block_post');
--
-- INSERT INTO group_permissions (group_id, permission_id)
-- SELECT
--     (SELECT id FROM groups WHERE name = 'moderators'),
--     id
-- FROM permissions
-- WHERE codename IN ('delete_post', 'view_post', 'block_post', 'hide_post');
--
--
-- INSERT INTO group_permissions (group_id, permission_id)
-- SELECT
--     (SELECT id FROM groups WHERE name = 'ordinary_users'),
--     id
-- FROM permissions
-- WHERE codename IN ('add_post', 'edit_post', 'delete_post', 'hide_post', 'view_post');