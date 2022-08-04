package common

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"time"
)

type Tx struct {
	tx *sql.Tx

	Debug  bool
	Driver string
	start  time.Time
}

func (t *Tx) fixQuery(query string) string {
	if t.Driver == "mysql" {
		return fixQuery(query)
	}
	return query
}

func (t *Tx) log(fname string, start time.Time, err error, tx bool, query string, args ...any) {
	if t.Debug {
		log(os.Stdout, fname, start, err, tx, query, args...)
	}
}

func (t *Tx) Commit() error {
	err := t.tx.Commit()
	t.log("Commit", t.start, err, true, "")
	return err
}

func (t *Tx) CurrentUnixTimestamp() int64 {
	return currentUnixTimestamp()
}

func (t *Tx) DeleteRowByID(ctx context.Context, id int64, row any) error {
	query := deleteRowByIDString(row)
	_, err := t.Exec(ctx, query, id)
	return err
}

func (t *Tx) Each(ctx context.Context, query string, callback func(ctx context.Context, rows *Rows) error, args ...any) error {
	if callback == nil {
		return fmt.Errorf("callback is not set")
	}
	rows, err := t.Query(ctx, query, args...)
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

func (t *Tx) EachPrepared(ctx context.Context, prep *Prepared, callback func(ctx context.Context, rows *Rows) error) error {
	return t.Each(ctx, prep.Query, callback, prep.Args...)
}

func (t *Tx) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := t.tx.ExecContext(ctx, t.fixQuery(query), args...)
	t.log("Exec", start, err, true, t.fixQuery(query), args...)
	return res, err
}

func (t *Tx) ExecPrepared(ctx context.Context, prep *Prepared) (sql.Result, error) {
	return t.Exec(ctx, prep.Query, prep.Args...)
}

func (t *Tx) InsertRow(ctx context.Context, row any) error {
	query, args := insertRowString(row)
	_, err := t.Exec(ctx, query, args...)
	return err
}

func (t *Tx) PrepareSQL(query string, args ...any) *Prepared {
	return prepareSQL(query, args...)
}

func (t *Tx) Query(ctx context.Context, query string, args ...any) (*Rows, error) {
	start := time.Now()
	rows, err := t.tx.QueryContext(ctx, t.fixQuery(query), args...)
	t.log("Query", start, err, true, t.fixQuery(query), args...)
	return &Rows{Rows: rows}, err
}

func (t *Tx) QueryPrepared(ctx context.Context, prep *Prepared) (*Rows, error) {
	return t.Query(ctx, prep.Query, prep.Args...)
}

func (t *Tx) QueryRow(ctx context.Context, query string, args ...any) *Row {
	start := time.Now()
	row := t.tx.QueryRowContext(ctx, t.fixQuery(query), args...)
	t.log("QueryRow", start, nil, true, t.fixQuery(query), args...)
	return &Row{Row: row}
}

func (t *Tx) QueryRowByID(ctx context.Context, id int64, row any) error {
	query := queryRowByIDString(row)
	return t.QueryRow(ctx, query, id).Scans(row)
}

func (t *Tx) QueryRowPrepared(ctx context.Context, prep *Prepared) *Row {
	return t.QueryRow(ctx, prep.Query, prep.Args...)
}

func (t *Tx) RowExists(ctx context.Context, id int64, row any) bool {
	var exists int
	query := rowExistsString(row)
	if err := t.QueryRow(ctx, query, id).Scan(&exists); err == nil && exists == 1 {
		return true
	}
	return false
}

func (t *Tx) Rollback() error {
	err := t.tx.Rollback()
	t.log("Rollback", t.start, err, true, "")
	return err
}

func (t *Tx) UpdateRow(ctx context.Context, row any) error {
	query, args := updateRowString(row)
	_, err := t.Exec(ctx, query, args...)
	return err
}

func (t *Tx) UpdateRowOnly(ctx context.Context, row any, only ...string) error {
	query, args := updateRowString(row, only...)
	_, err := t.Exec(ctx, query, args...)
	return err
}
