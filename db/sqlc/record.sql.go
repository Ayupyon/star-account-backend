// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: record.sql

package sqlc

import (
	"context"
	"time"

	"github.com/lib/pq"
)

const createRecord = `-- name: CreateRecord :one
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
) RETURNING id, name, type, date, amount, account_id, create_user_id, last_modified_user_id, create_time
`

type CreateRecordParams struct {
	Name         string    `json:"name"`
	Type         int32     `json:"type"`
	Date         time.Time `json:"date"`
	Amount       string    `json:"amount"`
	AccountID    int64     `json:"account_id"`
	CreateUserID int64     `json:"create_user_id"`
}

func (q *Queries) CreateRecord(ctx context.Context, arg CreateRecordParams) (Record, error) {
	row := q.db.QueryRowContext(ctx, createRecord,
		arg.Name,
		arg.Type,
		arg.Date,
		arg.Amount,
		arg.AccountID,
		arg.CreateUserID,
	)
	var i Record
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Date,
		&i.Amount,
		&i.AccountID,
		&i.CreateUserID,
		&i.LastModifiedUserID,
		&i.CreateTime,
	)
	return i, err
}

const deleteRecord = `-- name: DeleteRecord :exec
DELETE FROM records WHERE id=$1
`

func (q *Queries) DeleteRecord(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteRecord, id)
	return err
}

const deleteRecordsByAccountId = `-- name: DeleteRecordsByAccountId :exec
DELETE FROM records WHERE account_id=$1
`

func (q *Queries) DeleteRecordsByAccountId(ctx context.Context, accountID int64) error {
	_, err := q.db.ExecContext(ctx, deleteRecordsByAccountId, accountID)
	return err
}

const deleteRecordsByAccountIds = `-- name: DeleteRecordsByAccountIds :exec
DELETE FROM records WHERE account_id=ANY($1::bigint[])
`

func (q *Queries) DeleteRecordsByAccountIds(ctx context.Context, ids []int64) error {
	_, err := q.db.ExecContext(ctx, deleteRecordsByAccountIds, pq.Array(ids))
	return err
}

const getRecord = `-- name: GetRecord :one
SELECT id, name, type, date, amount, account_id, create_user_id, last_modified_user_id, create_time FROM records WHERE id=$1
`

func (q *Queries) GetRecord(ctx context.Context, id int64) (Record, error) {
	row := q.db.QueryRowContext(ctx, getRecord, id)
	var i Record
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Date,
		&i.Amount,
		&i.AccountID,
		&i.CreateUserID,
		&i.LastModifiedUserID,
		&i.CreateTime,
	)
	return i, err
}

const getRecordsAmountSumByAccountId = `-- name: GetRecordsAmountSumByAccountId :one
SELECT SUM(amount) FROM records WHERE account_id=$1
`

func (q *Queries) GetRecordsAmountSumByAccountId(ctx context.Context, accountID int64) (string, error) {
	row := q.db.QueryRowContext(ctx, getRecordsAmountSumByAccountId, accountID)
	var sum string
	err := row.Scan(&sum)
	return sum, err
}

const getRecordsByAccountId = `-- name: GetRecordsByAccountId :many
SELECT id, name, type, date, amount, account_id, create_user_id, last_modified_user_id, create_time FROM records
WHERE account_id=$1
OFFSET $2
LIMIT $3
`

type GetRecordsByAccountIdParams struct {
	AccountID int64 `json:"account_id"`
	Offset    int64 `json:"offset"`
	Limit     int64 `json:"limit"`
}

func (q *Queries) GetRecordsByAccountId(ctx context.Context, arg GetRecordsByAccountIdParams) ([]Record, error) {
	rows, err := q.db.QueryContext(ctx, getRecordsByAccountId, arg.AccountID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Record{}
	for rows.Next() {
		var i Record
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Date,
			&i.Amount,
			&i.AccountID,
			&i.CreateUserID,
			&i.LastModifiedUserID,
			&i.CreateTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecordsByAccountIdAndCreateUserId = `-- name: GetRecordsByAccountIdAndCreateUserId :many
SELECT id, name, type, date, amount, account_id, create_user_id, last_modified_user_id, create_time FROM records
WHERE account_id=$1 AND create_user_id=$2
OFFSET $3
LIMIT $4
`

type GetRecordsByAccountIdAndCreateUserIdParams struct {
	AccountID    int64 `json:"account_id"`
	CreateUserID int64 `json:"create_user_id"`
	Offset       int64 `json:"offset"`
	Limit        int64 `json:"limit"`
}

func (q *Queries) GetRecordsByAccountIdAndCreateUserId(ctx context.Context, arg GetRecordsByAccountIdAndCreateUserIdParams) ([]Record, error) {
	rows, err := q.db.QueryContext(ctx, getRecordsByAccountIdAndCreateUserId,
		arg.AccountID,
		arg.CreateUserID,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Record{}
	for rows.Next() {
		var i Record
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Date,
			&i.Amount,
			&i.AccountID,
			&i.CreateUserID,
			&i.LastModifiedUserID,
			&i.CreateTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecordsByAccountIdAndLastModifiedUserId = `-- name: GetRecordsByAccountIdAndLastModifiedUserId :many
SELECT id, name, type, date, amount, account_id, create_user_id, last_modified_user_id, create_time FROM records
WHERE account_id=$1 AND last_modified_user_id=$2
OFFSET $3
LIMIT $4
`

type GetRecordsByAccountIdAndLastModifiedUserIdParams struct {
	AccountID          int64 `json:"account_id"`
	LastModifiedUserID int64 `json:"last_modified_user_id"`
	Offset             int64 `json:"offset"`
	Limit              int64 `json:"limit"`
}

func (q *Queries) GetRecordsByAccountIdAndLastModifiedUserId(ctx context.Context, arg GetRecordsByAccountIdAndLastModifiedUserIdParams) ([]Record, error) {
	rows, err := q.db.QueryContext(ctx, getRecordsByAccountIdAndLastModifiedUserId,
		arg.AccountID,
		arg.LastModifiedUserID,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Record{}
	for rows.Next() {
		var i Record
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Date,
			&i.Amount,
			&i.AccountID,
			&i.CreateUserID,
			&i.LastModifiedUserID,
			&i.CreateTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecordsCountByAccountId = `-- name: GetRecordsCountByAccountId :one
SELECT COUNT(*) FROM records WHERE account_id=$1
`

func (q *Queries) GetRecordsCountByAccountId(ctx context.Context, accountID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getRecordsCountByAccountId, accountID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const updateRecord = `-- name: UpdateRecord :exec
UPDATE records
SET name=$2, type=$3, date=$4, amount=$5, last_modified_user_id=$6
WHERE id=$1
`

type UpdateRecordParams struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	Type               int32     `json:"type"`
	Date               time.Time `json:"date"`
	Amount             string    `json:"amount"`
	LastModifiedUserID int64     `json:"last_modified_user_id"`
}

func (q *Queries) UpdateRecord(ctx context.Context, arg UpdateRecordParams) error {
	_, err := q.db.ExecContext(ctx, updateRecord,
		arg.ID,
		arg.Name,
		arg.Type,
		arg.Date,
		arg.Amount,
		arg.LastModifiedUserID,
	)
	return err
}