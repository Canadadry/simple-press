package repository

import (
	"app/pkg/sqlutil"
)

type Repository struct {
	Db sqlutil.DBTX
}
