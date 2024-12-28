-- name: CreateSession :one
INSERT INTO sessions(
    session_data,
    user_agent,
    client_ip,
    is_blocked,
    expires_at,
    access_token,
    refresh_token,
    last_active

)VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
    RETURNING *;


-- name: BlockSession :one
UPDATE sessions
SET
    is_blocked = COALESCE(sqlc.narg('is_blocked'), is_blocked)
WHERE id = sqlc.arg('id')
    RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: UpdateSessionData :one
UPDATE sessions
SET
    access_token = COALESCE(sqlc.narg('access_token'), access_token),
    refresh_token = COALESCE(sqlc.narg('refresh_token'), refresh_token),
    session_data = COALESCE(sqlc.narg('session_data'), session_data),
    last_active = COALESCE(sqlc.narg('last_active'), last_active)
WHERE id = sqlc.arg('id')
    RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1;