package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"net/http"
)

const (
	templateAddName = "name"
)

type Template struct {
	Name string
}

func (t *Template) Bind(b validator.Binder) {
	b.RequiredStringVar(
		templateAddName,
		&t.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)
}

func ParseTemplateAdd(r *http.Request) (Template, validator.Errors, error) {
	parsed := Template{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return Template{}, validator.Errors{}, err
	}

	return parsed, errs, nil
}
