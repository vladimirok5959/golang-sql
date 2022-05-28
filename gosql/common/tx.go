package common

import (
	"context"
	"database/sql"
	"os"

	"time"
)

type Tx struct {
	tx *sql.Tx

	Debug  bool
	Driver string
	t      time.Time
}

func (db *Tx) fixQuery(query string) string {
	if db.Driver == "mysql" {
		return fixQuery(query)
	}
	return query
}

func (db *Tx) Commit() error {
	if db.Debug {
		err := db.tx.Commit()
		log(os.Stdout, "[func Commit]", db.t, err, true, "")
		return err
	}
	return db.tx.Commit()
}

func (db *Tx) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if db.Debug {
		t := time.Now()
		res, err := db.tx.ExecContext(ctx, db.fixQuery(query), args...)
		log(os.Stdout, "[func Exec]", t, err, true, db.fixQuery(query), args...)
		return res, err
	}
	return db.tx.ExecContext(ctx, db.fixQuery(query), args...)
}

func (db *Tx) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if db.Debug {
		t := time.Now()
		rows, err := db.tx.QueryContext(ctx, db.fixQuery(query), args...)
		log(os.Stdout, "[func Query]", t, err, true, db.fixQuery(query), args...)
		return rows, err
	}
	return db.tx.QueryContext(ctx, db.fixQuery(query), args...)
}

func (db *Tx) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	if db.Debug {
		t := time.Now()
		row := db.tx.QueryRowContext(ctx, db.fixQuery(query), args...)
		log(os.Stdout, "[func QueryRow]", t, nil, true, db.fixQuery(query), args...)
		return row
	}
	return db.tx.QueryRowContext(ctx, db.fixQuery(query), args...)
}

func (db *Tx) Rollback() error {
	if db.Debug {
		err := db.tx.Rollback()
		log(os.Stdout, "[func Rollback]", db.t, err, true, "")
		return err
	}
	return db.tx.Rollback()
}
