package common

import (
	"database/sql"
)

type Row struct {
	*sql.Row
}
