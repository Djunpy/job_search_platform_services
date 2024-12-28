-- name: CreateUser :one
INSERT INTO users (
    email,
    user_type,
    password,
    last_name,
    first_name,
    auth_source,
    sexy,
    last_token_update
) VALUES
    ($1, $2, $3, $4, $5, $6, $7, NOW())
    RETURNING *;

-- name: LastTokenUpdate :exec
UPDATE users
SET
    last_token_update = NOW()
WHERE id = $1;

-- name: UpdateUserById :one
UPDATE users
SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    verified_email = COALESCE(sqlc.narg('verified_email'), verified_email),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
    RETURNING *;

-- name: UpdateUserByEmail :one
UPDATE users
SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    verified_email = COALESCE(sqlc.narg('verified_email'), verified_email),
    updated_at = NOW()
WHERE email = sqlc.arg('email')
    RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserAndGroupsByEmail :one
SELECT
    u.id,
    u.email,
    u.verified_email,
    u.first_name,
    u.last_name,
    u.auth_source,
    u.date_joined,
    u.is_deleted,
    u.user_type,
    COALESCE(ARRAY_AGG(g.name), '{}') AS groups
FROM users u
         LEFT JOIN user_groups ug ON u.id = ug.user_id
         LEFT JOIN groups g ON ug.group_id = g.id
WHERE u.email = $1
GROUP BY
    u.id, u.email;


-- name: UserExists :one
SELECT EXISTS (
    SELECT * FROM users WHERE email = $1
);

-- name: DeleteUserByEmail :exec
DELETE FROM users WHERE email = $1;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1;


-- name: HideUserById :exec
UPDATE users
SET
    is_deleted = $2
WHERE id = $1;

-- name: HideUserByEmail :exec
UPDATE users
SET
    is_deleted = $2
WHERE email = $1;

-- name: ChangePassword :exec
UPDATE users
SET
    password = COALESCE(sqlc.narg('password'), password)
WHERE id = sqlc.arg('id');

-- name: NewUsersLast24H :one
SELECT COUNT(id) AS new_users_last_24h
FROM users
WHERE date_joined >= NOW() - INTERVAL '1 day';

-- name: GetAllUsersByRole :many
SELECT u.*
FROM users u
         JOIN user_groups ug ON u.id = ug.user_id
         JOIN groups g ON ug.group_id = g.id
WHERE g.name = $1
    LIMIT $2 OFFSET $3;

-- name: GetOrdinaryUsersCount :one
SELECT COUNT(u.id) AS ordinary_users_count
FROM users u
         JOIN user_groups ug ON u.id = ug.user_id
         JOIN groups g ON ug.group_id = g.id
WHERE g.name = $1;

-- name: FindUsers :many
SELECT u.*
FROM users u
         JOIN user_groups ug ON u.id = ug.user_id
         JOIN groups g ON ug.group_id = g.id
WHERE g.name = $1
  AND (u.email ILIKE '%' || $2 || '%')
    LIMIT $3 OFFSET $4;

-- name: GetAllUsersAndRoles :many
SELECT
    u.id,
    u.email,
    u.verified_email,
    u.auth_source,
    u.date_joined,
    u.is_deleted,
    u.user_type,
    COALESCE(ARRAY_AGG(g.name), '{}') AS groups
FROM users u
         LEFT JOIN user_groups ug ON u.id = ug.user_id
         LEFT JOIN groups g ON ug.group_id = g.id
WHERE g.name IN ('job_seekers', 'companies', 'premium_users')
GROUP BY
    u.id, u.email
ORDER BY
    u.id
    LIMIT $1 OFFSET $2;

-- name: GetAllUsers :many
SELECT
    u.id,
    u.email,
    u.verified_email,
    u.auth_source,
    u.date_joined,
    u.is_deleted,
    u.user_type
FROM users u
         LEFT JOIN user_groups ug ON u.id = ug.user_id
         LEFT JOIN groups g ON ug.group_id = g.id
WHERE g.name IN ('job_seekers', 'companies', 'premium_users')
GROUP BY
    u.id, u.email
ORDER BY
    u.id
    LIMIT $1 OFFSET $2;

-- name: SearchUsers :many
SELECT
    u.id,
    u.email,
    u.verified_email,
    u.auth_source,
    u.date_joined,
    u.is_deleted,
    u.user_type,
    COALESCE(ARRAY_AGG(g.name), '{}') AS groups
FROM
    users u
        LEFT JOIN
    user_groups ug ON u.id = ug.user_id
        LEFT JOIN
    groups g ON ug.group_id = g.id
WHERE
    (u.email ILIKE '%' || $1 || '%')
  AND g.name IN ('job_seekers', 'companies', 'premium_users')
GROUP BY
    u.id, u.email
ORDER BY
    u.id
    LIMIT $2 OFFSET $3;

-- name: CountPremiumUsers :one
SELECT COUNT(*) FROM groups g
                         JOIN user_groups ug ON g.id = ug.group_id
                         JOIN users u ON u.id = ug.user_id
WHERE g.name = 'premium_users';

-- name: CountUsersByGroup :one
SELECT COUNT(*) FROM groups g
                         JOIN user_groups ug ON g.id = ug.group_id
                         JOIN users u ON u.id = ug.user_id
WHERE g.name = $1;

-- name: CountUsersChurn30D :one
SELECT COUNT(*) FROM users
WHERE last_token_update < NOW() - INTERVAL '30 days' AND is_active = true;

-- name: UserGrowthPerYear :many
SELECT
    TO_CHAR(DATE_TRUNC('month', date_joined), 'MM') AS month,
    COUNT(*) AS user_count
FROM
    users
GROUP BY
    month
ORDER BY
    month;