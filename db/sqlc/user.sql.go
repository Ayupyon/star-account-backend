// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package sqlc

import (
	"context"

	"github.com/lib/pq"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    hashed_password
) VALUES (
    $1, $2, $3
) RETURNING id, name, email, hashed_password, create_time
`

type CreateUserParams struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Name, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.CreateTime,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, name, email, hashed_password, create_time FROM users WHERE id=$1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.CreateTime,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, hashed_password, create_time FROM users WHERE email=$1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.CreateTime,
	)
	return i, err
}

const getUsersByIds = `-- name: GetUsersByIds :many
SELECT id, name, email, hashed_password, create_time FROM users WHERE id=ANY($1::bigint[])
`

func (q *Queries) GetUsersByIds(ctx context.Context, ids []int64) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsersByIds, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.HashedPassword,
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

const getUsersByName = `-- name: GetUsersByName :many
SELECT id, name, email, hashed_password, create_time FROM users
WHERE name=$1
OFFSET $2
LIMIT $3
`

type GetUsersByNameParams struct {
	Name   string `json:"name"`
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
}

func (q *Queries) GetUsersByName(ctx context.Context, arg GetUsersByNameParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsersByName, arg.Name, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.HashedPassword,
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

const updateUserEmail = `-- name: UpdateUserEmail :exec
UPDATE users SET email=$2 WHERE id=$1
`

type UpdateUserEmailParams struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func (q *Queries) UpdateUserEmail(ctx context.Context, arg UpdateUserEmailParams) error {
	_, err := q.db.ExecContext(ctx, updateUserEmail, arg.ID, arg.Email)
	return err
}

const updateUserName = `-- name: UpdateUserName :exec
UPDATE users SET name=$2 WHERE id=$1
`

type UpdateUserNameParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateUserName(ctx context.Context, arg UpdateUserNameParams) error {
	_, err := q.db.ExecContext(ctx, updateUserName, arg.ID, arg.Name)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users SET hashed_password=$2 WHERE id=$1
`

type UpdateUserPasswordParams struct {
	ID             int64  `json:"id"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.ID, arg.HashedPassword)
	return err
}