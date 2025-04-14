package repository

import (
	"app/pkg/clock"
	"app/pkg/sqlutil"
)

type Repository struct {
	Db    sqlutil.DBTX
	Clock clock.Clock
}
