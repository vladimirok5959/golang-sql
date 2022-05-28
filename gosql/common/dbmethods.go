package common

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

type DBMethods struct {
	DB *sql.DB

	Debug  bool
	Driver string
}

func (db *DBMethods) fixQuery(query string) string {
	if db.Driver == "mysql" {
		return fixQuery(query)
	}
	return query
}

func (db *DBMethods) Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if db.Debug {
		t := time.Now()
		tx, err := db.DB.BeginTx(ctx, opts)
		log(os.Stdout, "Begin", t, err, true, "")
		return &Tx{tx, db.Debug, db.Driver, t}, err
	}

	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx, db.Debug, db.Driver, time.Now()}, err
}

func (db *DBMethods) Close() error {
	if db.Debug {
		t := time.Now()
		err := db.DB.Close()
		log(os.Stdout, "Close", t, err, false, "")
		return err
	}
	return db.DB.Close()
}

func (db *DBMethods) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if db.Debug {
		t := time.Now()
		res, err := db.DB.ExecContext(ctx, db.fixQuery(query), args...)
		log(os.Stdout, "Exec", t, err, false, db.fixQuery(query), args...)
		return res, err
	}
	return db.DB.ExecContext(ctx, db.fixQuery(query), args...)
}

func (db *DBMethods) Ping(ctx context.Context) error {
	if db.Debug {
		t := time.Now()
		err := db.DB.PingContext(ctx)
		log(os.Stdout, "Ping", t, err, false, "")
		return err
	}
	return db.DB.PingContext(ctx)
}

func (db *DBMethods) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	if db.Debug {
		t := time.Now()
		stm, err := db.DB.PrepareContext(ctx, db.fixQuery(query))
		log(os.Stdout, "Prepare", t, err, false, db.fixQuery(query))
		return stm, err
	}
	return db.DB.PrepareContext(ctx, db.fixQuery(query))
}

func (db *DBMethods) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if db.Debug {
		t := time.Now()
		rows, err := db.DB.QueryContext(ctx, db.fixQuery(query), args...)
		log(os.Stdout, "Query", t, err, false, db.fixQuery(query), args...)
		return rows, err
	}
	return db.DB.QueryContext(ctx, db.fixQuery(query), args...)
}

func (db *DBMethods) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	if db.Debug {
		t := time.Now()
		row := db.DB.QueryRowContext(ctx, db.fixQuery(query), args...)
		log(os.Stdout, "QueryRow", t, nil, false, db.fixQuery(query), args...)
		return row
	}
	return db.DB.QueryRowContext(ctx, db.fixQuery(query), args...)
}

func (db *DBMethods) SetConnMaxLifetime(d time.Duration) {
	db.DB.SetConnMaxLifetime(d)
}

func (db *DBMethods) SetMaxIdleConns(n int) {
	db.DB.SetMaxIdleConns(n)
}

func (db *DBMethods) SetMaxOpenConns(n int) {
	db.DB.SetMaxOpenConns(n)
}

func (db *DBMethods) Transaction(ctx context.Context, queries func(ctx context.Context, tx *Tx) error) error {
	if queries == nil {
		return fmt.Errorf("queries is not set for transaction")
	}
	tx, err := db.Begin(ctx, nil)
	if err != nil {
		return err
	}
	if err := queries(ctx, tx); err != nil {
		rerr := tx.Rollback()
		if rerr != nil {
			return fmt.Errorf(
				"%v: %v",
				rerr,
				err,
			)
		}
		return err
	}
	return tx.Commit()
}
