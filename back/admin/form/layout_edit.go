package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"net/http"
)

const (
	layoutEditName    = "name"
	layoutEditContent = "content"
)

type LayoutEdit struct {
	Name    string
	Content string
}

func (l *LayoutEdit) Bind(b validator.Binder) {
	b.RequiredStringVar(
		layoutEditName,
		&l.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)
	b.RequiredStringVar(
		layoutEditContent,
		&l.Content,
		validator.Length(1, maxContentLen),
	)
}

func ParseLayoutEdit(r *http.Request) (LayoutEdit, validator.Errors, error) {
	parsed := LayoutEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return LayoutEdit{}, validator.Errors{}, err
	}

	return parsed, errs, nil
}
