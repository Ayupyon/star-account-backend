-- name: CreateRecord :one
INSERT INTO records (
    name,
    type,
    date,
    amount,
    account_id,
    create_user_id,
    last_modified_user_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $6
) RETURNING *;

-- name: GetRecord :one
SELECT * FROM records WHERE id=$1;

-- name: GetRecordsByAccountId :many
SELECT * FROM records
WHERE account_id=$1
OFFSET $2
LIMIT $3;

-- name: GetRecordsCountByAccountId :one
SELECT COUNT(*) FROM records WHERE account_id=$1;

-- name: GetRecordsAmountSumByAccountId :one
SELECT SUM(amount) FROM records WHERE account_id=$1;

-- name: GetRecordsByAccountIdAndCreateUserId :many
SELECT * FROM records
WHERE account_id=$1 AND create_user_id=$2
OFFSET $3
LIMIT $4;

-- name: GetRecordsByAccountIdAndLastModifiedUserId :many
SELECT * FROM records
WHERE account_id=$1 AND last_modified_user_id=$2
OFFSET $3
LIMIT $4;

-- name: UpdateRecord :exec
UPDATE records
SET name=$2, type=$3, date=$4, amount=$5, last_modified_user_id=$6
WHERE id=$1;

-- name: DeleteRecord :exec
DELETE FROM records WHERE id=$1;

-- name: DeleteRecordsByAccountId :exec
DELETE FROM records WHERE account_id=$1;

-- name: DeleteRecordsByAccountIds :exec
DELETE FROM records WHERE account_id=ANY(sqlc.arg(ids)::bigint[]);