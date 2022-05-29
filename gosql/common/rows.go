package common

import (
	"database/sql"
)

type Rows struct {
	*sql.Rows
}

func (r *Rows) Scans(row any) error {
	return r.Rows.Scan(scans(row)...)
}
