package engine

import (
	"net/url"

	"github.com/vladimirok5959/golang-sql/gosql/common"
)

// ----------------------------------------------------------------------------

type mysql struct {
	*common.DBMethods
}

func NewMySQL(dbURL *url.URL, migrationsDir string, debug bool) (common.Engine, error) {
	db, err := common.OpenDB(dbURL, migrationsDir)
	if err != nil {
		return nil, err
	}

	return &mysql{
		DBMethods: &common.DBMethods{
			DB:     db,
			Debug:  debug,
			Driver: dbURL.Scheme,
		},
	}, nil
}

// ----------------------------------------------------------------------------

type postgresql struct {
	*common.DBMethods
}

func NewPostgreSQL(dbURL *url.URL, migrationsDir string, debug bool) (common.Engine, error) {
	db, err := common.OpenDB(dbURL, migrationsDir)
	if err != nil {
		return nil, err
	}

	return &postgresql{
		DBMethods: &common.DBMethods{
			DB:     db,
			Debug:  debug,
			Driver: dbURL.Scheme,
		},
	}, nil
}

// ----------------------------------------------------------------------------

type sqlite struct {
	*common.DBMethods
}

func NewSQLite(dbURL *url.URL, migrationsDir string, debug bool) (common.Engine, error) {
	db, err := common.OpenDB(dbURL, migrationsDir)
	if err != nil {
		return nil, err
	}

	return &sqlite{
		DBMethods: &common.DBMethods{
			DB:     db,
			Debug:  debug,
			Driver: dbURL.Scheme,
		},
	}, nil
}
