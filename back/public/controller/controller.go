package controller

import (
	"app/public/repository"
)

type Controller struct {
	Repository repository.Repository
}

func New(repo repository.Repository) (*Controller, error) {
	return &Controller{Repository: repo}, nil
}
