package common

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type DBMethods struct {
	DB *sql.DB

	Debug  bool
	Driver string
}

var rLogSpacesAll = regexp.MustCompile(`[\s\t]+`)
var rLogSpacesEnd = regexp.MustCompile(`[\s\t]+;$`)
var rSqlParam = regexp.MustCompile(`\$\d+`)

type queryFunc func(ctx context.Context, tx *sql.Tx) error

func (db *DBMethods) log(m string, s time.Time, e error, tx bool, query string, args ...any) {
	var tmsg string

	if tx {
		tmsg = " [TX]"
	}
	if m != "" {
		tmsg = tmsg + " " + m
	}

	qmsg := query
	if qmsg != "" {
		qmsg = strings.Trim(rLogSpacesAll.ReplaceAllString(qmsg, " "), " ")
		qmsg = rLogSpacesEnd.ReplaceAllString(qmsg, ";")
		qmsg = " " + qmsg
	}

	astr := " (empty)"
	if len(args) > 0 {
		astr = fmt.Sprintf(" (%v)", args)
	}

	estr := " (nil)"
	if e != nil {
		estr = " \033[0m\033[0;31m(" + e.Error() + ")"
	}

	color := "0;33"
	if tx {
		color = "1;33"
	}

	fmt.Fprintln(os.Stdout, "\033["+color+"m[SQL]"+tmsg+qmsg+astr+estr+fmt.Sprintf(" %.3f ms", time.Since(s).Seconds())+"\033[0m")
}

func (db *DBMethods) fixQuery(query string) string {
	if db.Driver == "mysql" {
		return rSqlParam.ReplaceAllString(query, "?")
	}
	return query
}

func (db *DBMethods) Begin(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if db.Debug {
		t := time.Now()
		tx, err := db.DB.BeginTx(ctx, opts)
		db.log("[func Begin]", t, err, true, "")
		return tx, err
	}
	return db.DB.BeginTx(ctx, opts)
}

func (db *DBMethods) Close() error {
	if db.Debug {
		t := time.Now()
		err := db.DB.Close()
		db.log("[func Close]", t, err, false, "")
		return err
	}
	return db.DB.Close()
}

func (db *DBMethods) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if db.Debug {
		t := time.Now()
		res, err := db.DB.ExecContext(ctx, db.fixQuery(query), args...)
		db.log("[func Exec]", t, err, false, db.fixQuery(query), args...)
		return res, err
	}
	return db.DB.ExecContext(ctx, db.fixQuery(query), args...)
}

func (db *DBMethods) Ping(ctx context.Context) error {
	if db.Debug {
		t := time.Now()
		err := db.DB.PingContext(ctx)
		db.log("[func Ping]", t, err, false, "")
		return err
	}
	return db.DB.PingContext(ctx)
}

func (db *DBMethods) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	if db.Debug {
		t := time.Now()
		stm, err := db.DB.PrepareContext(ctx, db.fixQuery(query))
		db.log("[func Prepare]", t, err, false, db.fixQuery(query))
		return stm, err
	}
	return db.DB.PrepareContext(ctx, db.fixQuery(query))
}

func (db *DBMethods) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if db.Debug {
		t := time.Now()
		rows, err := db.DB.QueryContext(ctx, db.fixQuery(query), args...)
		db.log("[func Query]", t, err, false, db.fixQuery(query), args...)
		return rows, err
	}
	return db.DB.QueryContext(ctx, db.fixQuery(query), args...)
}

func (db *DBMethods) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	if db.Debug {
		t := time.Now()
		row := db.DB.QueryRowContext(ctx, db.fixQuery(query), args...)
		db.log("[func QueryRow]", t, nil, false, db.fixQuery(query), args...)
		return row
	}
	return db.DB.QueryRowContext(ctx, db.fixQuery(query), args...)
}

func (db *DBMethods) Transaction(ctx context.Context, queries queryFunc) error {
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
