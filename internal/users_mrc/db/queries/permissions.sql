-- name: CreateGroup :one
INSERT INTO groups(id, name) VALUES ($1, $2)
    RETURNING *;

-- name: GetGroupByName :one
SELECT * FROM groups
WHERE name = $1;

-- name: UpdateGroup :one
UPDATE groups
SET
    name = COALESCE(sqlc.narg('name'), name)
WHERE id = $1
    RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;

-- name: AddUserToGroup :one
INSERT INTO user_groups(user_id, group_id)
VALUES ($1, $2)
    RETURNING *;

-- name: RemoveUserFromGroup :exec
DELETE FROM user_groups
WHERE user_id = $1 AND group_id = $2;

-- name: CreateUserGroup :one
INSERT INTO user_groups (user_id, group_id) VALUES ($1, $2)
    RETURNING *;

-- name: GetGroupsByUserId :many
SELECT
    g.name AS group_name,
    g.id AS group_id
FROM user_groups ug
         JOIN groups g ON g.id = ug.group_id
WHERE ug.user_id = $1;