-- name: CreateAccountAccessRule :one
INSERT INTO account_access_rules (
    user_id,
    account_id,
    role
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccountAccessRuleByUserIdAndAccountId :one
SELECT * FROM account_access_rules WHERE user_id=$1 AND account_id=$2;

-- name: GetUserIdsByAccountIdAndRole :many
SELECT user_id FROM account_access_rules
WHERE account_id=$1 AND role=$2
OFFSET $3
LIMIT $4;

-- name: GetUsersCountByAccountIdAndRole :one
SELECT COUNT(*) FROM account_access_rules WHERE account_id=$1 AND role=$2;

-- name: GetAccountIdsByUserIdAndRole :many
SELECT account_id FROM account_access_rules
WHERE user_id=$1 AND role=$2
OFFSET $3
LIMIT $4;

-- name: GetAccountsCountByUserIdAndRole :one
SELECT COUNT(*) FROM account_access_rules WHERE user_id=$1 AND role=$2;

-- name: GetAccountIdsByCreateUserIdForDelete :many
SELECT account_id FROM account_access_rules 
WHERE user_id=$1 AND role=2
FOR UPDATE;

-- name: DeleteAccountAccessRule :exec
DELETE FROM account_access_rules WHERE id=$1;

-- name: DeleteAccountAccessRulesByUserId :exec
DELETE FROM account_access_rules WHERE user_id=$1;

-- name: DeleteAccountAccessRulesByAccountId :exec
DELETE FROM account_access_rules WHERE account_id=$1;

-- name: DeleteAccountAccessRulesByAccountIds :exec
DELETE FROM account_access_rules WHERE account_id=ANY(sqlc.arg(ids)::bigint[]);

-- name: DeleteAccountAccessRuleByUserIdAndAccountId :exec
DELETE FROM account_access_rules WHERE user_id=$1 AND account_id=$2;