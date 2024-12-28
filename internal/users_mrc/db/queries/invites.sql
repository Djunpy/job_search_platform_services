-- name: CreateInvite :one
INSERT INTO invites(
    invite_code, created_by_user_id, expiration_date, group_id
) VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: UpdateInvite :one
UPDATE invites
SET
    used_by_user_id = COALESCE(sqlc.narg('used_by_user_id'), used_by_user_id),
    is_used = COALESCE(sqlc.narg('is_used'), is_used),
    updated_at = NOW()
WHERE invite_code = sqlc.arg('invite_code')
    RETURNING *;

-- name: DeleteInvite :exec
DELETE FROM invites
WHERE id = $1;

-- name: GetAllInvites :many
SELECT * FROM invites;

-- name: GetInvitesByUserId :many
SELECT * FROM invites
WHERE created_by_user_id = $1;

-- name: GetInviteByInviteCode :one
SELECT * FROM invites
WHERE invite_code = $1;