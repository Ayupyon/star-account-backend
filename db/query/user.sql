-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    hashed_password
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id=$1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=$1 LIMIT 1;

-- name: GetUsersByName :many
SELECT * FROM users
WHERE name=$1
OFFSET $2
LIMIT $3;

-- name: GetUsersByIds :many
SELECT * FROM users WHERE id=ANY(sqlc.arg(ids)::bigint[]);

-- name: UpdateUserName :exec
UPDATE users SET name=$2 WHERE id=$1;

-- name: UpdateUserEmail :exec
UPDATE users SET email=$2 WHERE id=$1;

-- name: UpdateUserPassword :exec
UPDATE users SET hashed_password=$2 WHERE id=$1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1;