-- name: GetUserPhoneByUserId :one
SELECT * FROM phones
WHERE user_id = $1;

-- name: CreateUserPhone :one
INSERT INTO phones (
    user_id,
    number,
    country_code
) VALUES ($1, $2, $3)
    RETURNING *;


-- name: UpdateUserPhone :one
UPDATE phones
SET
    number = coalesce(sqlc.narg('number'), number),
    country_code = coalesce(sqlc.narg('country_code'), country_code)
WHERE user_id = $1
    RETURNING *;

-- name: DeletePhoneByUserId :exec
DELETE FROM phones WHERE user_id = $1;