package common

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

type DBMethods struct {
	DB *sql.DB

	Debug  bool
	Driver string
}

type qFunc func(ctx context.Context, tx *sql.Tx) error

var r = regexp.MustCompile(`\$\d+`)

func (db *DBMethods) fixQuery(query string) string {
	if db.Driver == "mysql" {
		return r.ReplaceAllString(query, "?")
	}
	return query
}

func (db *DBMethods) Begin(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, opts)
}

func (db *DBMethods) Close() error {
	return db.DB.Close()
}

func (db *DBMethods) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.DB.ExecContext(ctx, db.fixQuery(query), args...)
}

func (db *DBMethods) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

func (db *DBMethods) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return db.DB.PrepareContext(ctx, db.fixQuery(query))
}

func (db *DBMethods) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, db.fixQuery(query), args...)
}

func (db *DBMethods) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return db.DB.QueryRowContext(ctx, db.fixQuery(query), args...)
}

func (db *DBMethods) Transaction(ctx context.Context, queries qFunc) error {
	if queries == nil {
		return fmt.Errorf("queries is not set for transaction")
	}
	tx, err := db.Begin(ctx, nil)
	if err != nil {
		return err
	}
	if err := queries(ctx, tx); err != nil {
		return errors.Wrap(err, tx.Rollback().Error())
	}
	return tx.Commit()
}
