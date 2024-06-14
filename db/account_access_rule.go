package db

import (
	"context"

	"github.com/timelyrain/star-account/db/sqlc"
)

type AccountAccessRule = sqlc.AccountAccessRule

func (db *DB) CreateAccountAccessRule(
	ctx context.Context,
	userId int64,
	accountId int64,
	role int32,
) (AccountAccessRule, error) {
	var res AccountAccessRule

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		arg := sqlc.CreateAccountAccessRuleParams{
			UserID:    userId,
			AccountID: accountId,
			Role:      role,
		}

		res, err = q.CreateAccountAccessRule(ctx, arg)
		return err
	})

	return res, err
}

func (db *DB) DeleteAccountAccessRule(ctx context.Context, id int64) error {
	return db.exec(ctx, func(q *sqlc.Queries) error {
		return q.DeleteAccountAccessRule(ctx, id)
	})
}

func (db *DB) GetAccountAccessRuleByUserIdAndAccountId(
	ctx context.Context,
	userId int64,
	accountId int64,
) (AccountAccessRule, error) {
	var res AccountAccessRule

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		arg := sqlc.GetAccountAccessRuleByUserIdAndAccountIdParams{
			UserID:    userId,
			AccountID: accountId,
		}

		res, err = q.GetAccountAccessRuleByUserIdAndAccountId(ctx, arg)
		return err
	})

	return res, err
}
