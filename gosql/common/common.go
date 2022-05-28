package common

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/url"

	"github.com/amacneil/dbmate/pkg/dbmate"
	_ "github.com/amacneil/dbmate/pkg/driver/mysql"
	_ "github.com/amacneil/dbmate/pkg/driver/postgres"
	_ "github.com/amacneil/dbmate/pkg/driver/sqlite"
	"golang.org/x/exp/slices"
)

type Engine interface {
	Close() error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Ping(context.Context) error
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row
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
