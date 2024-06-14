package db

import (
	"context"

	"github.com/timelyrain/star-account/db/sqlc"
)

type Account = sqlc.Account

/**
 * 1. 创建账户
 * 2. 创建对应于账户拥有者的权限信息
 */
func (db *DB) CreateAccount(ctx context.Context, name string, userId int64) (Account, error) {
	var res Account

	err := db.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		res, err = q.CreateAccount(ctx, name)
		if err != nil {
			return err
		}

		createAccountAccessRuleArg := sqlc.CreateAccountAccessRuleParams{
			UserID:    userId,
			AccountID: res.ID,
			Role:      2,
		}
		_, err = q.CreateAccountAccessRule(ctx, createAccountAccessRuleArg)
		return err
	})

	return res, err
}

/**
 * 1. 删除账单中的所有账单记录
 * 2. 删除账单中的所有权限信息
 * 3. 删除账单信息
 */
func (db *DB) DeleteAccount(ctx context.Context, id int64) error {
	return db.execTx(ctx, func(q *sqlc.Queries) error {
		err := q.DeleteRecordsByAccountId(ctx, id)
		if err != nil {
			return err
		}

		err = q.DeleteAccountAccessRulesByAccountId(ctx, id)
		if err != nil {
			return err
		}

		return q.DeleteAccount(ctx, id)
	})
}

func (db *DB) GetAccount(ctx context.Context, id int64) (Account, error) {
	var res Account

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		res, err = q.GetAccount(ctx, id)
		return err
	})

	return res, err
}

func (db *DB) GetAccountsCountByUserIdAndRole(ctx context.Context, userId int64, role int32) (int64, error) {
	var res int64

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		arg := sqlc.GetAccountsCountByUserIdAndRoleParams{
			UserID: userId,
			Role:   role,
		}

		res, err = q.GetAccountsCountByUserIdAndRole(ctx, arg)
		return err
	})
	return res, err
}

func (db *DB) GetAccountsByUserIdAndRole(
	ctx context.Context,
	userId int64,
	role int32,
	offset int64,
	limit int64,
) ([]Account, error) {
	var res []Account

	err := db.execTx(ctx, func(q *sqlc.Queries) error {
		var ids []int64
		var err error
		getIdsArg := sqlc.GetAccountIdsByUserIdAndRoleParams{
			UserID: userId,
			Role:   role,
			Offset: offset,
			Limit:  limit,
		}

		ids, err = q.GetAccountIdsByUserIdAndRole(ctx, getIdsArg)
		if err != nil {
			return err
		}

		res, err = q.GetAccountsByIds(ctx, ids)
		return err
	})

	return res, err
}

func (db *DB) UpdateAccountName(ctx context.Context, id int64, name string) error {
	return db.exec(ctx, func(q *sqlc.Queries) error {
		arg := sqlc.UpdateAccountNameParams{
			ID:   id,
			Name: name,
		}
		return q.UpdateAccountName(ctx, arg)
	})
}
