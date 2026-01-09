package controller

import (
	"app/admin/repository"
	"app/pkg/clock"
)

type Controller struct {
	Repository repository.Repository
	Clock      clock.Clock
}

func New(repo repository.Repository, c clock.Clock) (*Controller, error) {
	return &Controller{Repository: repo, Clock: c}, nil
}
