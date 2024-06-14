-- name: CreateAccount :one
INSERT INTO accounts (
    name
) VALUES (
    $1
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE id=$1;

-- name: GetAccountsByIds :many
SELECT * FROM accounts WHERE id=ANY(sqlc.arg(ids)::bigint[]);

-- name: UpdateAccountName :exec
UPDATE accounts SET name=$2 WHERE id=$1;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id=$1;

-- name: DeleteAccountsByIds :exec
DELETE FROM accounts WHERE id=ANY(sqlc.arg(ids)::bigint[]);