package controller

import (
	"app/admin/repository"
	"app/admin/translation"
	"app/pkg/clock"
	"app/pkg/i18n"
	"fmt"
)

type Controller struct {
	Repository repository.Repository
	tr         i18n.Translator
	Clock      clock.Clock
}

func New(repo repository.Repository, c clock.Clock) (*Controller, error) {
	tr, err := translation.GetTranslator()
	if err != nil {
		return nil, fmt.Errorf("can't load translations : %w", err)
	}
	return &Controller{Repository: repo, Clock: c, tr: tr}, nil
}
