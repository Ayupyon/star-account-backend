package db

import (
	"context"
	"time"

	"github.com/timelyrain/star-account/db/sqlc"
)

type Record = sqlc.Record

func (db *DB) CreateRecord(
	ctx context.Context,
	name string,
	recordType int32,
	date time.Time,
	amount string,
	accountId int64,
	createUserId int64,
) (Record, error) {
	var res Record

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		arg := sqlc.CreateRecordParams{
			Name:         name,
			Type:         recordType,
			Date:         date,
			Amount:       amount,
			AccountID:    accountId,
			CreateUserID: createUserId,
		}
		res, err = q.CreateRecord(ctx, arg)
		return err
	})

	return res, err
}

func (db *DB) DeleteRecord(ctx context.Context, id int64) error {
	return db.exec(ctx, func(q *sqlc.Queries) error {
		return q.DeleteRecord(ctx, id)
	})
}

func (db *DB) GetRecord(ctx context.Context, id int64) (Record, error) {
	var res Record

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		res, err = q.GetRecord(ctx, id)
		return err
	})

	return res, err
}

func (db *DB) GetRecordsByAccountId(
	ctx context.Context,
	accountId int64,
	offset int64,
	limit int64,
) ([]Record, error) {
	var res []Record

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		arg := sqlc.GetRecordsByAccountIdParams{
			AccountID: accountId,
			Offset:    offset,
			Limit:     limit,
		}

		res, err = q.GetRecordsByAccountId(ctx, arg)
		return err
	})

	return res, err
}

func (db *DB) GetRecordsCountByAccountId(ctx context.Context, accountId int64) (int64, error) {
	var res int64

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		res, err = q.GetRecordsCountByAccountId(ctx, accountId)
		return err
	})

	return res, err
}

func (db *DB) GetRecordsAmountSumByAccountId(
	ctx context.Context,
	accountId int64,
) (string, error) {
	var res string

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		res, err = q.GetRecordsAmountSumByAccountId(ctx, accountId)
		return err
	})

	return res, err
}

func (db *DB) UpdateRecord(
	ctx context.Context,
	id int64,
	name string,
	recordType int32,
	date time.Time,
	amount string,
	lastModifiedUserId int64,
) error {
	return db.exec(ctx, func(q *sqlc.Queries) error {
		arg := sqlc.UpdateRecordParams{
			ID:                 id,
			Name:               name,
			Type:               recordType,
			Date:               date,
			Amount:             amount,
			LastModifiedUserID: lastModifiedUserId,
		}
		return q.UpdateRecord(ctx, arg)
	})
}
