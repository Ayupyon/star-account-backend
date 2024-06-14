package db

import (
	"context"

	"github.com/timelyrain/star-account/db/sqlc"
)

type User = sqlc.User

func (db *DB) CreateUser(
	ctx context.Context,
	name string,
	email string,
	hashedPassword string,
) (User, error) {
	var res User

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		arg := sqlc.CreateUserParams{
			Name:           name,
			Email:          email,
			HashedPassword: hashedPassword,
		}

		res, err = q.CreateUser(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})

	return res, err
}

/**
 * 1. 找到用户拥有的所有账单的id
 * 2. 按照账单id删除所有的账单记录和账单权限信息
 * 3. 删除用户id关联的所有账单权限信息(该用户是管理者,但不是拥有者)
 * 4. 删除用户拥有的所有账单
 * 5. 删除该用户的信息
 */
func (db *DB) DeleteUser(ctx context.Context, id int64) error {
	return db.execTx(ctx, func(q *sqlc.Queries) error {
		var accountIds []int64
		var err error

		accountIds, err = q.GetAccountIdsByCreateUserIdForDelete(ctx, id)
		if err != nil {
			return err
		}

		err = q.DeleteRecordsByAccountIds(ctx, accountIds)
		if err != nil {
			return err
		}

		err = q.DeleteAccountAccessRulesByAccountIds(ctx, accountIds)
		if err != nil {
			return err
		}

		err = q.DeleteAccountAccessRulesByUserId(ctx, id)
		if err != nil {
			return err
		}

		err = q.DeleteAccountsByIds(ctx, accountIds)
		if err != nil {
			return err
		}

		return q.DeleteUser(ctx, id)
	})
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var res User

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		res, err = q.GetUserByEmail(ctx, email)
		return err
	})

	return res, err
}

func (db *DB) GetUser(ctx context.Context, id int64) (User, error) {
	var res User

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		res, err = q.GetUser(ctx, id)
		return err
	})

	return res, err
}

func (db *DB) GetUsersByAccountIdAndRole(
	ctx context.Context,
	accountId int64,
	role int32,
	offset int64,
	limit int64,
) ([]User, error) {
	var res []User

	err := db.execTx(ctx, func(q *sqlc.Queries) error {
		var ids []int64
		var err error
		getIdArg := sqlc.GetUserIdsByAccountIdAndRoleParams{
			AccountID: accountId,
			Role:      role,
			Offset:    offset,
			Limit:     limit,
		}

		ids, err = q.GetUserIdsByAccountIdAndRole(ctx, getIdArg)
		if err != nil {
			return err
		}

		res, err = q.GetUsersByIds(ctx, ids)
		return err
	})

	return res, err
}

func (db *DB) GetUsersByName(
	ctx context.Context,
	name string,
	offset int64,
	limit int64,
) ([]User, error) {
	var res []User

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		arg := sqlc.GetUsersByNameParams{
			Name:   name,
			Offset: offset,
			Limit:  limit,
		}

		res, err = q.GetUsersByName(ctx, arg)
		return err
	})

	return res, err
}

func (db *DB) UpdateUserEmail(ctx context.Context, id int64, email string) error {
	return db.exec(ctx, func(q *sqlc.Queries) error {
		arg := sqlc.UpdateUserEmailParams{
			ID:    id,
			Email: email,
		}
		return q.UpdateUserEmail(ctx, arg)
	})
}

func (db *DB) UpdateUserName(ctx context.Context, id int64, name string) error {
	return db.exec(ctx, func(q *sqlc.Queries) error {
		arg := sqlc.UpdateUserNameParams{
			ID:   id,
			Name: name,
		}
		return q.UpdateUserName(ctx, arg)
	})
}

func (db *DB) UpdateUserPassword(ctx context.Context, id int64, hashedPassword string) error {
	return db.exec(ctx, func(q *sqlc.Queries) error {
		arg := sqlc.UpdateUserPasswordParams{
			ID:             id,
			HashedPassword: hashedPassword,
		}
		return q.UpdateUserPassword(ctx, arg)
	})
}

func (db *DB) GetUsersCountByAccountIdAndRole(
	ctx context.Context,
	accountId int64,
	role int32,
) (int64, error) {
	var res int64

	err := db.exec(ctx, func(q *sqlc.Queries) error {
		var err error
		arg := sqlc.GetUsersCountByAccountIdAndRoleParams{
			AccountID: accountId,
			Role:      role,
		}

		res, err = q.GetUsersCountByAccountIdAndRole(ctx, arg)
		return err
	})

	return res, err
}
