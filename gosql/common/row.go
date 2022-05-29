package common

import (
	"database/sql"
)

type Row struct {
	*sql.Row
}

func (r *Row) Scans(row any) error {
	return r.Row.Scan(scans(row)...)
}
