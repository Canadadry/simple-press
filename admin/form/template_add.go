package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"net/http"
	"strings"
)

const (
	templateAddName = "name"
)

type Template struct {
	Name string
}

type TemplateError struct {
	Name string
}

func (te TemplateError) HasError() bool {
	return te.Name != ""
}

func (t *Template) Bind(b validator.Binder) {
	b.RequiredStringVar(
		templateAddName,
		&t.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)
}

func ParseTemplateAdd(r *http.Request) (Template, TemplateError, error) {
	parsed := Template{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return Template{}, TemplateError{}, err
	}

	resultErr := TemplateError{
		Name: strings.Join(errs.Errors[templateAddName], ", "),
	}

	return parsed, resultErr, nil
}
