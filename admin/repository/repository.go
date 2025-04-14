package repository

import (
	"app/pkg/clock"
	"app/pkg/sqlutil"
)

type Repository struct {
	db    sqlutil.DBTX
	clock clock.Clock
}

func New(db sqlutil.DBTX, clock clock.Clock) *Repository {
	return &Repository{
		db:    db,
		clock: clock,
	}
}
