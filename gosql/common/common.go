package common

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"os"
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

func log(w io.Writer, fname string, start time.Time, err error, tx bool, query string, args ...any) string {
	var values []string

	bold := "0"
	color := "33"

	// Transaction or not
	if tx {
		bold = "1"
		values = append(values, "[TX]")
	}

	// Function name
	if fname != "" {
		values = append(values, fname)
	}

	// SQL query
	if query != "" {
		values = append(values, rLogSpacesEnd.ReplaceAllString(
			strings.Trim(rLogSpacesAll.ReplaceAllString(query, " "), " "), ";",
		))
	}

	// Params
	if len(args) > 0 {
		values = append(values, fmt.Sprintf("(%v)", args))
	} else {
		values = append(values, "(empty)")
	}

	// Error
	if err != nil {
		color = "31"
		values = append(values, "("+err.Error()+")")
	} else {
		values = append(values, "(nil)")
	}

	// Execute time with close color symbols
	values = append(values, fmt.Sprintf("%.3f ms\033[0m", time.Since(start).Seconds()))

	// Prepend start caption with colors
	values = append([]string{"\033[" + bold + ";" + color + "m[SQL]"}, values...)

	res := fmt.Sprintln(strings.Join(values, " "))
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

func OpenDB(databaseURL *url.URL, migrationsDir string, debug bool) (*sql.DB, error) {
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

	var db *sql.DB

	if debug {
		t := time.Now()
		db, err = driver.Open()
		log(os.Stdout, "[func Open]", t, err, false, "")
		if err != nil {
			return nil, fmt.Errorf("DB open error: %w", err)
		}
	} else {
		db, err = driver.Open()
		if err != nil {
			return nil, fmt.Errorf("DB open error: %w", err)
		}
	}

	return db, nil
}
