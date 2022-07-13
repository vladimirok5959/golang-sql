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

func (d *DBMethods) fixQuery(query string) string {
	if d.Driver == "mysql" {
		return fixQuery(query)
	}
	return query
}

func (d *DBMethods) log(fname string, start time.Time, err error, tx bool, query string, args ...any) {
	if d.Debug {
		log(os.Stdout, fname, start, err, tx, query, args...)
	}
}

func (d *DBMethods) Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	start := time.Now()
	tx, err := d.DB.BeginTx(ctx, opts)
	d.log("Begin", start, err, true, "")
	return &Tx{tx, d.Debug, d.Driver, start}, err
}

func (d *DBMethods) Close() error {
	start := time.Now()
	err := d.DB.Close()
	d.log("Close", start, err, false, "")
	return err
}

func (d *DBMethods) DeleteRowByID(ctx context.Context, id int64, row any) error {
	query := deleteRowByIDString(row)
	_, err := d.Exec(ctx, query, id)
	return err
}

func (d *DBMethods) Each(ctx context.Context, query string, callback func(ctx context.Context, rows *Rows) error, args ...any) error {
	if callback == nil {
		return fmt.Errorf("callback is not set")
	}
	rows, err := d.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := callback(ctx, rows); err != nil {
				return err
			}
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func (d *DBMethods) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := d.DB.ExecContext(ctx, d.fixQuery(query), args...)
	d.log("Exec", start, err, false, d.fixQuery(query), args...)
	return res, err
}

func (d *DBMethods) Ping(ctx context.Context) error {
	start := time.Now()
	err := d.DB.PingContext(ctx)
	d.log("Ping", start, err, false, "")
	return err
}

func (d *DBMethods) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	start := time.Now()
	stm, err := d.DB.PrepareContext(ctx, d.fixQuery(query))
	d.log("Prepare", start, err, false, d.fixQuery(query))
	return stm, err
}

func (d *DBMethods) Query(ctx context.Context, query string, args ...any) (*Rows, error) {
	start := time.Now()
	rows, err := d.DB.QueryContext(ctx, d.fixQuery(query), args...)
	d.log("Query", start, err, false, d.fixQuery(query), args...)
	return &Rows{Rows: rows}, err
}

func (d *DBMethods) QueryRow(ctx context.Context, query string, args ...any) *Row {
	start := time.Now()
	row := d.DB.QueryRowContext(ctx, d.fixQuery(query), args...)
	d.log("QueryRow", start, nil, false, d.fixQuery(query), args...)
	return &Row{Row: row}
}

func (d *DBMethods) QueryRowByID(ctx context.Context, id int64, row any) error {
	query := queryRowByIDString(row)
	return d.QueryRow(ctx, query, id).Scans(row)
}

func (d *DBMethods) RowExists(ctx context.Context, id int64, row any) bool {
	var exists int
	query := rowExistsString(row)
	if err := d.QueryRow(ctx, query, id).Scan(&exists); err == nil && exists == 1 {
		return true
	}
	return false
}

func (d *DBMethods) SetConnMaxLifetime(t time.Duration) {
	start := time.Now()
	d.DB.SetConnMaxLifetime(t)
	d.log("SetConnMaxLifetime", start, nil, false, "")
}

func (d *DBMethods) SetMaxIdleConns(n int) {
	start := time.Now()
	d.DB.SetMaxIdleConns(n)
	d.log("SetMaxIdleConns", start, nil, false, "")
}

func (d *DBMethods) SetMaxOpenConns(n int) {
	start := time.Now()
	d.DB.SetMaxOpenConns(n)
	d.log("SetMaxOpenConns", start, nil, false, "")
}

func (d *DBMethods) Transaction(ctx context.Context, callback func(ctx context.Context, tx *Tx) error) error {
	if callback == nil {
		return fmt.Errorf("callback is not set")
	}
	tx, err := d.Begin(ctx, nil)
	if err != nil {
		return err
	}
	if err := callback(ctx, tx); err != nil {
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
