package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"net/http"
	"strings"
)

const (
	layoutAddName = "name"
)

type Layout struct {
	Name string
}

type LayoutError struct {
	Name string
}

func (le LayoutError) HasError() bool {
	return le.Name != ""
}

func (l *Layout) Bind(b validator.Binder) {
	b.RequiredStringVar(
		layoutAddName,
		&l.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)
}

func ParseLayoutAdd(r *http.Request) (Layout, LayoutError, error) {
	parsed := Layout{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return Layout{}, LayoutError{}, err
	}

	resultErr := LayoutError{
		Name: strings.Join(errs.Errors[layoutAddName], ", "),
	}

	return parsed, resultErr, nil
}
