package gosql

import (
	"fmt"

	"github.com/vladimirok5959/golang-sql/gosql/common"
	"github.com/vladimirok5959/golang-sql/gosql/engine"
)

func Open(dbURL, migrationsDir string) (common.Engine, error) {
	databaseURL, err := common.ParseUrl(dbURL)
	if err != nil {
		return nil, err
	}

	switch databaseURL.Scheme {
	case "mysql":
		return engine.NewMySQL(databaseURL, migrationsDir)
	case "postgres", "postgresql":
		return engine.NewPostgreSQL(databaseURL, migrationsDir)
	case "sqlite", "sqlite3":
		return engine.NewSQLite(databaseURL, migrationsDir)
	default:
		return nil, fmt.Errorf("DB open error")
	}
}
