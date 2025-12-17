package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"net/http"
)

const (
	layoutAddName = "name"
)

type Layout struct {
	Name string
}

func (l *Layout) Bind(b validator.Binder) {
	b.RequiredStringVar(
		layoutAddName,
		&l.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)
}

func ParseLayoutAdd(r *http.Request) (Layout, validator.Errors, error) {
	parsed := Layout{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return Layout{}, validator.Errors{}, err
	}

	return parsed, errs, nil
}
