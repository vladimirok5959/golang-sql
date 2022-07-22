package common

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
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
	CurrentUnixTimestamp() int64
	DeleteRowByID(ctx context.Context, id int64, row any) error
	Each(ctx context.Context, query string, logic func(ctx context.Context, rows *Rows) error, args ...any) error
	EachPrepared(ctx context.Context, prep *Prepared, logic func(ctx context.Context, rows *Rows) error) error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	ExecPrepared(ctx context.Context, prep *Prepared) (sql.Result, error)
	Ping(context.Context) error
	Prepare(ctx context.Context, query string) (*sql.Stmt, error)
	Query(ctx context.Context, query string, args ...any) (*Rows, error)
	QueryPrepared(ctx context.Context, prep *Prepared) (*Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *Row
	QueryRowByID(ctx context.Context, id int64, row any) error
	QueryRowPrepared(ctx context.Context, prep *Prepared) *Row
	RowExists(ctx context.Context, id int64, row any) bool
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Transaction(ctx context.Context, queries func(ctx context.Context, tx *Tx) error) error
}

var rSqlParam = regexp.MustCompile(`\$\d+`)
var rLogSpacesAll = regexp.MustCompile(`[\s\t]+`)
var rLogSpacesEnd = regexp.MustCompile(`[\s\t]+;$`)

func currentUnixTimestamp() int64 {
	return time.Now().UTC().Unix()
}

func deleteRowByIDString(row any) string {
	v := reflect.ValueOf(row).Elem()
	t := v.Type()
	var table string
	for i := 0; i < t.NumField(); i++ {
		if table == "" {
			if tag := t.Field(i).Tag.Get("table"); tag != "" {
				table = tag
			}
		}
	}
	return `DELETE FROM ` + table + ` WHERE id = $1`
}

func fixQuery(query string) string {
	return rSqlParam.ReplaceAllString(query, "?")
}

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
		values = append(values, "[func "+fname+"]")
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

func queryRowByIDString(row any) string {
	v := reflect.ValueOf(row).Elem()
	t := v.Type()
	var table string
	fields := []string{}
	for i := 0; i < t.NumField(); i++ {
		if table == "" {
			if tag := t.Field(i).Tag.Get("table"); tag != "" {
				table = tag
			}
		}
		if tag := t.Field(i).Tag.Get("field"); tag != "" {
			fields = append(fields, tag)
		}
	}
	return `SELECT ` + strings.Join(fields, ", ") + ` FROM ` + table + ` WHERE id = $1 LIMIT 1`
}

func rowExistsString(row any) string {
	v := reflect.ValueOf(row).Elem()
	t := v.Type()
	var table string
	for i := 0; i < t.NumField(); i++ {
		if table == "" {
			if tag := t.Field(i).Tag.Get("table"); tag != "" {
				table = tag
			}
		}
	}
	return `SELECT 1 FROM ` + table + ` WHERE id = $1 LIMIT 1`
}

func scans(row any) []any {
	v := reflect.ValueOf(row).Elem()
	res := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		res[i] = v.Field(i).Addr().Interface()
	}
	return res
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

func OpenDB(databaseURL *url.URL, migrationsDir string, skipMigration bool, debug bool) (*sql.DB, error) {
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

	if !skipMigration {
		if err := mate.CreateAndMigrate(); err != nil {
			return nil, fmt.Errorf("DB migration error: %w", err)
		}
	}

	var db *sql.DB
	start := time.Now()
	db, err = driver.Open()
	if debug {
		log(os.Stdout, "Open", start, err, false, "")
	}
	if err != nil {
		return nil, fmt.Errorf("DB open error: %w", err)
	}

	return db, nil
}

func PrepareSQL(query string, args ...any) *Prepared {
	return &Prepared{query, args}
}
