package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/timelyrain/star-account/db/sqlc"
)

type DB struct {
	conn *sql.DB
}

func NewDB(conn *sql.DB) *DB {
	return &DB{conn: conn}
}

func (db *DB) exec(ctx context.Context, fn func(*sqlc.Queries) error) error {
	q := sqlc.New(db.conn)
	return fn(q)
}

func (db *DB) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return errors.Join(err, rbErr)
		}

		return err
	}

	return tx.Commit()
}
