package gosql

import (
	"fmt"

	"github.com/vladimirok5959/golang-sql/gosql/common"
	"github.com/vladimirok5959/golang-sql/gosql/engine"
)

type Row = common.Row

type Rows = common.Rows

type Tx = common.Tx

func Open(dbURL, migrationsDir string, debug bool) (common.Engine, error) {
	databaseURL, err := common.ParseUrl(dbURL)
	if err != nil {
		return nil, err
	}

	switch databaseURL.Scheme {
	case "mysql":
		return engine.NewMySQL(databaseURL, migrationsDir, debug)
	case "postgres", "postgresql":
		return engine.NewPostgreSQL(databaseURL, migrationsDir, debug)
	case "sqlite", "sqlite3":
		return engine.NewSQLite(databaseURL, migrationsDir, debug)
	default:
		return nil, fmt.Errorf("DB open error")
	}
}
