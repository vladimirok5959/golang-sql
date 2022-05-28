package common

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/amacneil/dbmate/pkg/dbmate"
	_ "github.com/amacneil/dbmate/pkg/driver/mysql"
	_ "github.com/amacneil/dbmate/pkg/driver/postgres"
	_ "github.com/amacneil/dbmate/pkg/driver/sqlite"
	"golang.org/x/exp/slices"
)

type Engine interface {
	Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error)
	Close() error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Ping(context.Context) error
	Prepare(ctx context.Context, query string) (*sql.Stmt, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row
	Transaction(ctx context.Context, queries func(ctx context.Context, tx *Tx) error) error
}

var rLogSpacesAll = regexp.MustCompile(`[\s\t]+`)
var rLogSpacesEnd = regexp.MustCompile(`[\s\t]+;$`)
var rSqlParam = regexp.MustCompile(`\$\d+`)

func log(w io.Writer, m string, s time.Time, e error, tx bool, query string, args ...any) string {
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

	bold := "0"
	color := "33"

	estr := " (nil)"
	if e != nil {
		color = "31"
		estr = " (" + e.Error() + ")"
	}

	if tx {
		bold = "1"
	}

	res := fmt.Sprintln("\033[" + bold + ";" + color + "m[SQL]" + tmsg + qmsg + astr + estr + fmt.Sprintf(" %.3f ms", time.Since(s).Seconds()) + "\033[0m")
	fmt.Fprint(w, res)
	return res
}

func fixQuery(query string) string {
	return rSqlParam.ReplaceAllString(query, "?")
}

func ParseUrl(dbURL string) (*url.URL, error) {
	databaseURL, err := url.Parse(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %w", err)
	}

	if databaseURL.Scheme == "" {
		return nil, fmt.Errorf("protocol scheme is not defined")
	}

	protocols := []string{"mysql", "postgres", "postgresql", "sqlite", "sqlite3"}
	if !slices.Contains(protocols, databaseURL.Scheme) {
		return nil, fmt.Errorf("unsupported protocol scheme: %s", databaseURL.Scheme)
	}

	return databaseURL, nil
}

func OpenDB(databaseURL *url.URL, migrationsDir string) (*sql.DB, error) {
	mate := dbmate.New(databaseURL)

	mate.AutoDumpSchema = false
	mate.Log = io.Discard
	if migrationsDir != "" {
		mate.MigrationsDir = migrationsDir
	}

	driver, err := mate.GetDriver()
	if err != nil {
		return nil, fmt.Errorf("DB get driver error: %w", err)
	}

	if err := mate.CreateAndMigrate(); err != nil {
		return nil, fmt.Errorf("DB migration error: %w", err)
	}

	db, err := driver.Open()
	if err != nil {
		return nil, fmt.Errorf("DB open error: %w", err)
	}

	return db, nil
}
